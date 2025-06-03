package main

import (
	"blive-vup-layer/config"
	"blive-vup-layer/dao"
	"blive-vup-layer/llm"
	"blive-vup-layer/speechrecognition"
	"blive-vup-layer/tts"
	"blive-vup-layer/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/golang-lru/v2/expirable"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/vtb-link/bianka/basic"
	"github.com/vtb-link/bianka/live"
	"github.com/vtb-link/bianka/proto"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/exp/slog"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

const (
	FansMedalName = "巫女酱" // 粉丝牌名称

	LlmReplyFansMedalLevel     = 10 // 可以触发大模型响应的最小粉丝牌等级
	RoomEnterTTSFansMedalLevel = 15 // 可以触发进入直播间TTS提示的最小粉丝牌等级

	MessageExpiration  = 15 * time.Minute // 历史消息过期时间
	GiftComboDuration  = 4 * time.Second  // 礼物连击时间，连击结束后会合并播放TTS
	LlmHistoryDuration = 10 * time.Minute // 大模型使用历史弹幕去理解上下文的时间范围
	TimeIdleDuration   = 10 * time.Minute // 直播间进入空闲的时间

	DisableLlmByUserCountDuration = 1 * time.Minute // 统计间隔时间内用户数量，用于触发暂停大模型
	DisableLlmByUserCount         = 5               // 触发暂停大模型的用户数量

	LlmReplyLimitDuration = 5 * time.Minute // 大模型最大回复数量的统计时间
	LlmReplyLimitCount    = 10              // 大模型统计窗口内最大的回复数量

	ProbabilityLlmTriggerDuration    = 5 * time.Minute // 概率触发大模型回复的统计时间
	ProbabilityLlmTriggerLevel1      = 0.0             // 100%触发
	ProbabilityLlmTriggerLevel1Count = 0               // 统计人数为0
	ProbabilityLlmTriggerLevel2      = 0.3             // 70%触发
	ProbabilityLlmTriggerLevel2Count = 10              // 统计人数为[1, 10]
	ProbabilityLlmTriggerLevel3      = 0.7             // 30%触发
)

type Service struct {
	cfg        *config.Config
	liveClient *live.Client

	LLM *llm.LLM
	TTS *tts.TTS
	Sr  *speechrecognition.SpeechRecognition
	Dao *dao.Dao

	slog *slog.Logger

	app        *App
	appContext context.Context

	livingCfg   LiveConfig
	startResp   *live.AppStartResponse
	tk          *time.Ticker
	wcs         *basic.WsClient
	connContext context.Context
	connCancel  context.CancelFunc

	historyMsgLru               *expirable.LRU[string, *ChatMessage]
	llmReplyLru                 *expirable.LRU[string, struct{}]
	probabilityLlmTriggerRandom *rand.Rand
	isLiving                    bool
	ttsQueue                    *tts.TTSQueue
	ttsCh                       <-chan *tts.TaskResult
	lastEnterUser               *UserData
	roomIdleTimer               *time.Timer

	areaName  string // 直播分区
	roomTitle string // 直播间标题
}

func (s *Service) Init(app *App) {
	s.app = app

	const (
		ConfigProdFilePath = "./etc/config.toml"
		ConfigDevFilePath  = "./etc/config-dev.toml"
	)

	var configFilePath string
	envInfo := app.App.Environment()
	if envInfo.Debug {
		configFilePath = ConfigDevFilePath
	} else {
		configFilePath = ConfigProdFilePath
	}

	var err error
	s.cfg, err = config.ParseConfig(configFilePath)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
		return
	}

	s.TTS, err = tts.NewTTS(s.cfg.AliyunTTS)
	if err != nil {
		log.Fatalf("tts.NewTTS err: %v", err)
		return
	}
	s.ttsQueue = tts.NewTTSQueue(s.TTS)

	s.Sr, err = speechrecognition.NewSpeechRecognition(s.cfg.SpeechRecognition)
	if err != nil {
		log.Fatalf("speechrecognition.NewSpeechRecognition err: %v", err)
		return
	}

	s.Dao, err = dao.NewDao(s.cfg.DbPath)
	if err != nil {
		log.Fatalf("dao.NewDao err: %v", err)
		return
	}

	s.liveClient = live.NewClient(live.NewConfig(s.cfg.BiliBili.AccessKey, s.cfg.BiliBili.SecretKey, s.cfg.BiliBili.AppId))
	s.LLM = llm.NewLLM(s.cfg.LLM)

	s.ttsCh = s.ttsQueue.ListenResult()
	util.RunGr(func() {
		for r := range s.ttsCh {
			if err := r.Err; err != nil {
				s.writeResultError(ResultTypeTTS, CodeInternalError, err.Error())
				continue
			}
			s.writeResultOK(ResultTypeTTS, map[string]interface{}{
				"audio_file_path": r.Fname,
			})
			s.roomIdleTimer.Reset(TimeIdleDuration)
		}
	})
}

func (s *Service) OnStartup(ctx context.Context, options application.ServiceOptions) error {
	s.appContext = ctx
	return nil
}

func NewService(logWriter io.Writer) *Service {
	return &Service{
		slog: slog.New(slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})),

		livingCfg: LiveConfig{
			DisableTTS:          false,
			DisableLlm:          false,
			DisableWelcomeLimit: false,
			DisableIdleTTS:      false,
		},

		roomIdleTimer: time.NewTimer(TimeIdleDuration),

		historyMsgLru:               expirable.NewLRU[string, *ChatMessage](512, nil, MessageExpiration),
		llmReplyLru:                 expirable.NewLRU[string, struct{}](LlmReplyLimitCount, nil, LlmReplyLimitDuration),
		probabilityLlmTriggerRandom: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer util.Recover()

	fileName := r.URL.Path
	filePath := path.Join(config.ResultFilePath, fileName)
	f, err := os.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Errorf("fileName: %s, filePath: %s not found", fileName, filePath)
		return
	}

	w.Write(f)
}

type GiftWithTimer struct {
	Uname    string
	GiftName string
	GiftNum  int32
	Timer    *time.Timer
}

type LiveConfig struct {
	DisableTTS          bool `json:"disable_tts"`
	DisableLlm          bool `json:"disable_llm"`
	DisableWelcomeLimit bool `json:"disable_welcome_limit"`
	DisableIdleTTS      bool `json:"disable_idle_tts"`
}

type ChatMessage struct {
	User      *UserData
	Message   string
	Timestamp time.Time
}

func (s *Service) InitConn(initData *InitRequestData) *Result {
	defer util.Recover()

	s.livingCfg = initData.Config
	s.writeResultOK(ResultTypeConfig, s.livingCfg)

	s.init(initData.Code)

	return BuildResultOk(nil)
}

func (s *Service) SetConfig(configData LiveConfig) *Result {
	defer util.Recover()

	s.livingCfg = configData
	return BuildResultOk(s.livingCfg)
}

func (s *Service) pushTTS(params *tts.NewTaskParams, force bool) {
	if (!s.isLiving && !force) || s.livingCfg.DisableTTS {
		return
	}
	log.Infof("pushTTS text: %s", params.Text)
	if err := s.ttsQueue.Push(params); err != nil {
		s.writeResultError(ResultTypeTTS, CodeInternalError, err.Error())
	}
}

func (s *Service) init(code string) {
	if s.startResp != nil {
		s.writeResultOK(ResultTypeRoom, &RoomData{
			RoomID: s.startResp.AnchorInfo.RoomID,
			Uname:  s.startResp.AnchorInfo.Uname,
			UFace:  s.startResp.AnchorInfo.UFace,
		})
		return
	}

	log.Infof("init code: %s", code)
	startResp, err := s.liveClient.AppStart(code)
	if err != nil {
		s.writeResultError(ResultTypeRoom, http.StatusBadRequest, err.Error())
		return
	}
	s.startResp = startResp

	s.connContext, s.connCancel = context.WithCancel(s.appContext)

	s.tk = time.NewTicker(time.Second * 20)
	util.RunGr(func() {
		for {
			select {
			case <-s.appContext.Done():
				return
			case <-s.tk.C:
				// 心跳
				if err := s.liveClient.AppHeartbeat(startResp.GameInfo.GameID); err != nil {
					log.Errorf("Heartbeat fail, err: %v", err)
					s.connCancel()
					s.StopConn()
					util.RunGr(func() {
						s.init(code)
					})
					return
				}
			}
		}
	})

	// close 事件处理
	onCloseHandle := func(wcs *basic.WsClient, startResp basic.StartResp, closeType int) {
		defer util.Recover()

		// 注册关闭回调
		log.Infof("WebsocketClient onClose, startResp: %v", startResp)

		// 注意检查关闭类型, 避免无限重连
		if closeType == basic.CloseActively || closeType == basic.CloseReceivedShutdownMessage || closeType == basic.CloseAuthFailed {
			log.Infof("WebsocketClient exit")
			return
		}

		// 对于可能的情况下重新连接
		// 注意: 在某些场景下 startResp 会变化, 需要重新获取
		// 此外, 一但 AppHeartbeat 失败, 会导致 startResp.GameInfo.GameID 变化, 需要重新获取
		err := wcs.Reconnection(startResp)
		if err != nil {
			log.Errorf("Reconnection fail, err: %v", err)
			s.writeResultError(ResultTypeRoom, CodeInternalError, err.Error())
			s.connCancel()
			s.StopConn()
			util.RunGr(func() {
				s.init(code)
			})
			return
		}
	}

	s.lastEnterUser = nil

	util.RunGr(func() {
		for range s.roomIdleTimer.C {
			if s.livingCfg.DisableIdleTTS {
				s.roomIdleTimer.Reset(TimeIdleDuration)
				continue
			}
			text := util.GetRandomStr(RandomIdleReply)
			s.pushTTS(&tts.NewTaskParams{Text: text}, false)
			s.roomIdleTimer.Reset(TimeIdleDuration)
		}
	})

	s.isLiving = true

	giftTimerMap := make(map[string]*GiftWithTimer)
	var giftTimerMapMutex sync.RWMutex

	// 消息处理 Handle
	dispatcherHandleMap := basic.DispatcherHandleMap{
		proto.OperationMessage: func(_ *basic.WsClient, msg *proto.Message) error {
			defer util.Recover()

			// 单条消息raw
			log.Infof(string(msg.Payload()))

			// 自动解析
			_, data, err := proto.AutomaticParsingMessageCommand(msg.Payload())
			if err != nil {
				log.Errorf("proto.AutomaticParsingMessageCommand err: %v", err)
				return err
			}

			// Switch cmd
			switch d := data.(type) {
			case *proto.CmdDanmuData:
				{
					if _, ok := danmuGiftMap[d.Msg]; ok {
						break
					}
					u := UserData{
						OpenID:                 d.OpenID,
						Uname:                  d.Uname,
						UFace:                  d.UFace,
						FansMedalLevel:         d.FansMedalLevel,
						FansMedalName:          d.FansMedalName,
						FansMedalWearingStatus: d.FansMedalWearingStatus,
						GuardLevel:             d.GuardLevel,
					}
					danmuData := &DanmuData{
						UserData:    u,
						Msg:         d.Msg,
						MsgID:       d.MsgID,
						Timestamp:   d.Timestamp,
						EmojiImgUrl: d.EmojiImgUrl,
						DmType:      d.DmType,
					}
					s.writeResultOK(ResultTypeDanmu, danmuData)

					util.RunGr(func() {
						s.setUser(u)
					})

					s.historyMsgLru.Add(d.MsgID, &ChatMessage{
						User:      &u,
						Message:   danmuData.Msg,
						Timestamp: time.Now(),
					})

					pitchRate := 0
					//if !livingCfg.DisableLlm {
					//	pitchRate = -100
					//}
					s.pushTTS(&tts.NewTaskParams{
						Text:      fmt.Sprintf("%s说：%s", d.Uname, d.Msg),
						PitchRate: pitchRate,
					}, false)

					if danmuData.FansMedalWearingStatus &&
						danmuData.FansMedalName == FansMedalName &&
						danmuData.FansMedalLevel >= LlmReplyFansMedalLevel { // 带10级粉丝牌
						s.startLlmReply(false)
					} else if danmuData.GuardLevel > 0 || // 舰长
						(danmuData.Uname == "巫女酱子" || danmuData.Uname == "青云-_-z") {
						s.startLlmReply(true)
					}

					break
				}
			case *proto.CmdSuperChatData:
				{
					u := UserData{
						OpenID:                 d.OpenID,
						Uname:                  d.Uname,
						UFace:                  d.Uface,
						FansMedalLevel:         d.FansMedalLevel,
						FansMedalName:          d.FansMedalName,
						FansMedalWearingStatus: d.FansMedalWearingStatus,
						GuardLevel:             d.GuardLevel,
					}
					scData := &SuperChatData{
						UserData:  u,
						Msg:       d.Message,
						MsgID:     d.MsgID,
						MessageID: d.MessageID,
						Rmb:       float64(d.Rmb),
						Timestamp: d.Timestamp,
						StartTime: d.StartTime,
						EndTime:   d.EndTime,
					}
					s.writeResultOK(ResultTypeSuperChat, scData)

					util.RunGr(func() {
						s.setUser(u)
					})

					s.historyMsgLru.Add(d.MsgID, &ChatMessage{
						User:      &u,
						Message:   scData.Msg,
						Timestamp: time.Now(),
					})
					s.pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("谢谢%s酱的醒目留言：%s", d.Uname, d.Message),
					}, false)
					s.startLlmReply(true)
					break
				}
			case *proto.CmdSendGiftData:
				{
					u := UserData{
						OpenID:                 d.OpenID,
						Uname:                  d.Uname,
						UFace:                  d.Uface,
						FansMedalLevel:         d.FansMedalLevel,
						FansMedalName:          d.FansMedalName,
						FansMedalWearingStatus: d.FansMedalWearingStatus,
						GuardLevel:             d.GuardLevel,
					}
					s.writeResultOK(ResultTypeGift, &GiftData{
						UserData:  u,
						GiftID:    d.GiftID,
						GiftName:  d.GiftName,
						GiftNum:   d.GiftNum,
						Rmb:       float64(d.Price) / 1000,
						Paid:      d.Paid,
						Timestamp: d.Timestamp,
						MsgID:     d.MsgID,
						GiftIcon:  d.GiftIcon,
						ComboGift: d.ComboGift,
						ComboInfo: &GiftDataComboInfo{
							ComboBaseNum: d.ComboInfo.ComboBaseNum,
							ComboCount:   d.ComboInfo.ComboCount,
							ComboID:      d.ComboInfo.ComboID,
							ComboTimeout: d.ComboInfo.ComboTimeout,
						},
					})

					util.RunGr(func() {
						s.setUser(u)
					})

					key := fmt.Sprintf("%s-%d", d.OpenID, d.GiftID)

					giftTimerMapMutex.RLock()
					gt, ok := giftTimerMap[key]
					giftTimerMapMutex.RUnlock()
					if ok {
						atomic.AddInt32(&gt.GiftNum, int32(d.GiftNum))
						gt.Timer.Reset(GiftComboDuration)
						break
					}

					gt = &GiftWithTimer{
						Uname:    d.Uname,
						GiftNum:  int32(d.GiftNum),
						GiftName: d.GiftName,
						Timer:    time.NewTimer(GiftComboDuration),
					}

					giftTimerMapMutex.Lock()
					giftTimerMap[key] = gt
					giftTimerMapMutex.Unlock()
					util.RunGr(func() {
						defer gt.Timer.Stop()
						<-gt.Timer.C

						giftTimerMapMutex.Lock()
						delete(giftTimerMap, key)
						giftTimerMapMutex.Unlock()

						giftNum := atomic.LoadInt32(&gt.GiftNum)
						s.pushTTS(&tts.NewTaskParams{
							Text: fmt.Sprintf("谢谢%s酱赠送的%d个%s 么么哒", gt.Uname, giftNum, gt.GiftName),
						}, false)
					})
					break
				}
			case *proto.CmdGuardData:
				{
					u := UserData{
						OpenID:                 d.UserInfo.OpenID,
						Uname:                  d.UserInfo.Uname,
						UFace:                  d.UserInfo.Uface,
						FansMedalLevel:         d.FansMedalLevel,
						FansMedalName:          d.FansMedalName,
						FansMedalWearingStatus: d.FansMedalWearingStatus,
						GuardLevel:             d.GuardLevel,
					}
					guardName := getGuardLevelName(d.GuardLevel)
					s.writeResultOK(ResultTypeGuard, &GuardData{
						UserData:   u,
						Rmb:        float64(d.Price) / 1000,
						GuardLevel: d.GuardLevel,
						GuardNum:   d.GuardNum,
						GuardUnit:  d.GuardUnit,
						GuardName:  guardName,
						Timestamp:  d.Timestamp,
						MsgID:      d.MsgID,
					})
					util.RunGr(func() {
						s.setUser(u)
					})
					s.pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("谢谢%s酱赠送的%d个%s%s，么么哒", d.UserInfo.Uname, d.GuardNum, d.GuardUnit, guardName),
					}, false)
					break
				}
			case *proto.CmdLiveStartData:
				{
					s.areaName = d.AreaName
					s.roomTitle = d.Title
					s.pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("主人开始直播啦，弹幕姬启动！直播分区为%s，直播间标题为%s", d.AreaName, d.Title),
					}, true)
					s.isLiving = true
					break
				}
			case *proto.CmdLiveEndData:
				{
					s.pushTTS(&tts.NewTaskParams{
						Text: "主人直播结束啦，今天辛苦了！",
					}, true)
					s.isLiving = false
					break
				}
			case *proto.CmdLiveRoomEnterData:
				{
					u := UserData{
						OpenID: d.OpenID,
						Uname:  d.Uname,
						UFace:  d.Uface,
					}
					s.writeResultOK(ResultTypeEnterRoom, &RoomEnterData{
						UserData:  u,
						MsgID:     uuid.NewV4().String(),
						Timestamp: d.Timestamp,
					})

					s.lastEnterUser = &u

					util.RunGr(func() {
						u, err := s.Dao.GetUser(context.Background(), u.OpenID)
						if err != nil {
							log.Errorf("GetUser open_id: %s err: %v", u.OpenID, err)
							return
						}

						if u == nil {
							return
						}

						if s.livingCfg.DisableWelcomeLimit ||
							(u.FansMedalWearingStatus && u.FansMedalLevel >= RoomEnterTTSFansMedalLevel) ||
							u.GuardLevel > 0 {

							name := d.Uname
							if u.GuardLevel > 0 {
								guardName := getGuardLevelName(u.GuardLevel)
								name = guardName + name
							}

							s.pushTTS(&tts.NewTaskParams{
								Text: fmt.Sprintf("欢迎%s酱来到直播间", name),
							}, false)
						}
					})

					break
				}
			case *proto.CmdRoomChangeData:
				{
					s.areaName = d.AreaName
					s.roomTitle = d.Title
					s.pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("直播间信息发生变更，直播分区为%s，直播间标题为%s", d.AreaName, d.Title),
					}, true)
					break
				}
			case *proto.CmdInteractWordData:
				{
					s.writeResultOK(ResultTypeInteractWord, &InteractWordData{
						MsgID:     d.MsgID,
						OpenID:    d.OpenID,
						RoomID:    d.RoomID,
						Timestamp: d.Timestamp,
						Uname:     d.Uname,
					})
					s.pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("谢谢%s酱关注直播间，么么哒", d.Uname),
					}, true)
					break
				}
			case *proto.CmdWarningData:
				{
					s.pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("直播间收到超管警告，警告信息为：%s", d.Msg),
					}, true)
				}
			case map[string]interface{}:
				{
					dataJson, err := json.Marshal(d)
					if err != nil {
						break
					}
					log.Infof("unknown data: %s", string(dataJson))
				}
			default:
				{
					break
				}
			}

			return nil
		},
	}

	s.wcs, err = basic.StartWebsocket(
		startResp,
		dispatcherHandleMap,
		onCloseHandle,
		s.slog,
	)
	if err != nil {
		log.Errorf("basic.StartWebsocket err: %v", err)
		s.writeResultError(ResultTypeRoom, CodeInternalError, err.Error())
		return
	}

	log.Infof("room_info: %v", startResp.AnchorInfo)
	s.writeResultOK(ResultTypeRoom, &RoomData{
		RoomID: startResp.AnchorInfo.RoomID,
		Uname:  startResp.AnchorInfo.Uname,
		UFace:  startResp.AnchorInfo.UFace,
	})
}

func (s *Service) startLlmReply(force bool) {
	if !s.isLiving || s.livingCfg.DisableLlm {
		return
	}

	var msgs []*ChatMessage
	userMap := map[string]struct{}{}
	probabilityLlmTriggerCounter := -1 // 当前尝试触发的用户不算，所以初始值为-1
	for _, msg := range s.historyMsgLru.Values() {
		if msg == nil {
			continue
		}
		if time.Since(msg.Timestamp) <= LlmHistoryDuration {
			msgs = append(msgs, msg)
		}
		if time.Since(msg.Timestamp) <= DisableLlmByUserCountDuration {
			userMap[msg.User.OpenID] = struct{}{}
		}
		if time.Since(msg.Timestamp) <= ProbabilityLlmTriggerDuration {
			probabilityLlmTriggerCounter++
		}
	}

	if len(msgs) == 0 {
		return
	}

	if !force {
		llmReplyLruLen := s.llmReplyLru.Len()
		if llmReplyLruLen >= LlmReplyLimitCount {
			log.Infof("disable llm by reply count: %d", llmReplyLruLen)
			return
		}

		if len(userMap) >= DisableLlmByUserCount {
			log.Infof("disable llm by user count: %d", len(userMap))
			return
		}

		currentMsg := msgs[len(msgs)-1]
		if util.IsRepeatedChar(currentMsg.Message) {
			log.Infof("disable llm by repeated msg: %s", currentMsg.Message)
			return
		}

		var probability float64
		if probabilityLlmTriggerCounter > ProbabilityLlmTriggerLevel2Count {
			probability = ProbabilityLlmTriggerLevel3
		} else if probabilityLlmTriggerCounter > ProbabilityLlmTriggerLevel1Count {
			probability = ProbabilityLlmTriggerLevel2
		} else {
			probability = ProbabilityLlmTriggerLevel1
		}

		r := s.probabilityLlmTriggerRandom.Float64()
		log.Infof("r: %.2f, probability: %.2f\n", r, probability)
		if r <= probability {
			log.Infof("disable llm by probability: %.2f, counter: %d, compare: %.2f", r, probabilityLlmTriggerCounter, probability)
			return
		}
	}

	util.RunGr(func() {
		llmMsgs := make([]*llm.ChatMessage, len(msgs))
		var lastMsg *ChatMessage
		for i, msg := range msgs {
			llmMsgs[i] = &llm.ChatMessage{
				User:    msg.User.Uname,
				Message: msg.Message,
			}
			lastMsg = msg
		}

		var extraInfo []string
		if s.areaName != "" {
			extraInfo = append(extraInfo, fmt.Sprintf("当前直播分区：%s", s.areaName))
		}
		if s.roomTitle != "" {
			extraInfo = append(extraInfo, fmt.Sprintf("当前直播间标题：%s", s.roomTitle))
		}
		llmRes, err := s.LLM.ChatWithLLM(context.Background(), &llm.ChatWithLLMParams{
			ExtraInfos: extraInfo,
			Messages:   llmMsgs,
		})
		if err != nil {
			log.Errorf("ChatWithLLM err: %v", err)
			return
		}
		msgId := uuid.NewV4().String()
		s.writeResultOK(ResultTypeLLM, LLMResult{
			UserData:    *lastMsg.User,
			MsgID:       msgId,
			UserMessage: lastMsg.Message,
			LLMResult:   llmRes.Content,
		})
		s.historyMsgLru.Add(msgId, &ChatMessage{
			User: &UserData{
				Uname: "巫女酱子的助手",
			},
			Message:   llmRes.Content,
			Timestamp: time.Now(),
		})
		s.llmReplyLru.Add(msgId, struct{}{})
		s.pushTTS(&tts.NewTaskParams{
			Text: llmRes.Content,
		}, false)
	})
}

func (s *Service) OnShutdown() error {
	s.StopConn()
	return nil
}

func (s *Service) StopConn() {
	defer util.Recover()

	if s.startResp != nil {
		s.liveClient.AppEnd(s.startResp.GameInfo.GameID)
	}
	if s.wcs != nil {
		s.wcs.Close()
		s.wcs = nil
	}
	s.startResp = nil
	s.connContext = nil
	if s.connCancel != nil {
		s.connCancel()
		s.connCancel = nil
	}
	if s.tk != nil {
		s.tk.Stop()
		s.tk = nil
	}
	s.lastEnterUser = nil
}

func (s *Service) setUser(userData UserData) {
	err := s.Dao.CreateOrUpdateUser(context.Background(), &dao.User{
		OpenID:                 userData.OpenID,
		FansMedalWearingStatus: userData.FansMedalWearingStatus,
		FansMedalLevel:         userData.FansMedalLevel,
		GuardLevel:             userData.GuardLevel,
	})
	if err != nil {
		log.Errorf("CreateOrUpdateUser open_id: %s, err: %v", userData.OpenID, err)
	}
}

func getGuardLevelName(guardLevel int) string {
	guardName, ok := GuardLevelMap[guardLevel]
	if !ok {
		guardName = "舰长"
	}
	return guardName
}
