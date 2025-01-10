package main

import (
	"blive-vup-layer/config"
	"blive-vup-layer/dao"
	"blive-vup-layer/llm"
	"blive-vup-layer/tts"
	"context"
	"fmt"
	"github.com/hashicorp/golang-lru/v2/expirable"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/vtb-link/bianka/basic"
	"github.com/vtb-link/bianka/live"
	"github.com/vtb-link/bianka/proto"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/exp/slog"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	FansMedalName = "巫女酱" // 粉丝牌名称

	LlmReplyFansMedalLevel     = 10 // 可以触发大模型响应的最小粉丝牌等级
	RoomEnterTTSFansMedalLevel = 15 // 可以触发进入直播间TTS提示的最小粉丝牌等级

	MessageExpiration     = 15 * time.Minute // 历史消息过期时间
	GiftComboDuration     = 4 * time.Second  // 礼物连击时间，连击结束后会合并播放TTS
	LlmHistoryDuration    = 10 * time.Minute // 大模型使用历史弹幕去理解上下文的时间范围
	LastEnterUserDuration = 10 * time.Minute // 最后一个进入直播间用户将会播放TTS的等待时间

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

type App struct {
	cfg        *config.Config
	liveClient *live.Client

	LLM *llm.LLM
	TTS *tts.TTS
	Dao *dao.Dao

	slog *slog.Logger

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
	isLlmProcessing             bool
	isLiving                    bool
	ttsQueue                    *tts.TTSQueue
	lastEnterUser               *UserData
	lastEnterUserTimer          *time.Timer
}

func (a *App) startup(ctx context.Context) {
	a.appContext = ctx

	const (
		ConfigProdFilePath = "./etc/config.toml"
		ConfigDevFilePath  = "./etc/config-dev.toml"
	)

	var configFilePath string
	envInfo := runtime.Environment(ctx)
	if envInfo.BuildType == "production" {
		configFilePath = ConfigProdFilePath
	} else {
		configFilePath = ConfigDevFilePath
	}

	var err error
	a.cfg, err = config.ParseConfig(configFilePath)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
		return
	}

	a.TTS, err = tts.NewTTS(a.cfg.AliyunTTS)
	if err != nil {
		log.Fatalf("tts.NewTTS err: %v", err)
		return
	}

	a.Dao, err = dao.NewDao(a.cfg.DbPath)
	if err != nil {
		log.Fatalf("dao.NewDao err: %v", err)
		return
	}

	a.liveClient = live.NewClient(live.NewConfig(a.cfg.BiliBili.AccessKey, a.cfg.BiliBili.SecretKey, a.cfg.BiliBili.AppId))
	a.LLM = llm.NewLLM(a.cfg.QianFan)
}

func NewApp(logWriter io.Writer) *App {
	return &App{
		slog: slog.New(slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: slog.LevelInfo})),

		livingCfg: LiveConfig{
			DisableLlm: false,
		},
	}
}

type GiftWithTimer struct {
	Uname    string
	GiftName string
	GiftNum  int32
	Timer    *time.Timer
}

type LiveConfig struct {
	DisableLlm bool `json:"disable_llm"`
}

type ChatMessage struct {
	OpenId    string
	User      string
	Message   string
	Timestamp time.Time
}

func (a *App) InitConn(initData *InitRequestData) *Result {
	a.livingCfg = initData.Config
	a.writeResultOK(ResultTypeConfig, a.livingCfg)

	a.init(initData.Code)

	return BuildResultOk(nil)
}

func (a *App) SetConfig(configData LiveConfig) *Result {
	a.livingCfg = configData
	return BuildResultOk(a.livingCfg)
}

func (a *App) pushTTS(params *tts.NewTaskParams, force bool) {
	if !a.isLiving && !force {
		return
	}
	if err := a.ttsQueue.Push(params); err != nil {
		a.writeResultError(ResultTypeTTS, CodeInternalError, err.Error())
	}
}

func (a *App) init(code string) {
	if a.startResp != nil {
		a.writeResultOK(ResultTypeRoom, &RoomData{
			RoomID: a.startResp.AnchorInfo.RoomID,
			Uname:  a.startResp.AnchorInfo.Uname,
			UFace:  a.startResp.AnchorInfo.UFace,
		})
		return
	}

	log.Infof("init code: %s", code)
	startResp, err := a.liveClient.AppStart(code)
	if err != nil {
		a.writeResultError(ResultTypeRoom, http.StatusBadRequest, err.Error())
		return
	}
	a.startResp = startResp

	a.connContext, a.connCancel = context.WithCancel(a.appContext)

	a.tk = time.NewTicker(time.Second * 20)
	go func() {
		for {
			select {
			case <-a.appContext.Done():
				return
			case <-a.tk.C:
				// 心跳
				if err := a.liveClient.AppHeartbeat(startResp.GameInfo.GameID); err != nil {
					log.Errorf("Heartbeat fail, err: %v", err)
					a.connCancel()
					a.StopConn()
					go a.init(code)
					return
				}
			}
		}
	}()

	// close 事件处理
	onCloseHandle := func(wcs *basic.WsClient, startResp basic.StartResp, closeType int) {
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
			a.writeResultError(ResultTypeRoom, CodeInternalError, err.Error())
			a.connCancel()
			a.StopConn()
			go a.init(code)
			return
		}
	}

	a.lastEnterUser = nil
	a.lastEnterUserTimer = time.NewTimer(LastEnterUserDuration)

	a.ttsQueue = tts.NewTTSQueue(a.TTS)
	ttsCh := a.ttsQueue.ListenResult()
	go func() {
		for r := range ttsCh {
			if err := r.Err; err != nil {
				a.writeResultError(ResultTypeTTS, CodeInternalError, err.Error())
				continue
			}
			a.writeResultOK(ResultTypeTTS, map[string]interface{}{
				"audio_file_path": r.Fname,
			})
			a.lastEnterUserTimer.Reset(LastEnterUserDuration)
		}
	}()
	go func() {
		for r := range ttsCh {
			if err := r.Err; err != nil {
				a.writeResultError(ResultTypeTTS, CodeInternalError, err.Error())
				continue
			}
			a.writeResultOK(ResultTypeTTS, map[string]interface{}{
				"audio_file_path": r.Fname,
			})
			a.lastEnterUserTimer.Reset(LastEnterUserDuration)
		}
	}()
	pushTTS := func(params *tts.NewTaskParams, force bool) {
		if !a.isLiving && !force {
			return
		}
		if err := a.ttsQueue.Push(params); err != nil {
			a.writeResultError(ResultTypeTTS, CodeInternalError, err.Error())
		}
	}

	go func() {
		for range a.lastEnterUserTimer.C {
			if a.lastEnterUser == nil {
				a.lastEnterUserTimer.Reset(LastEnterUserDuration)
				continue
			}
			pushTTS(&tts.NewTaskParams{
				Text: fmt.Sprintf("欢迎%s酱来到直播间", a.lastEnterUser.Uname),
			}, false)
			a.lastEnterUserTimer.Reset(LastEnterUserDuration)
		}
	}()

	a.historyMsgLru = expirable.NewLRU[string, *ChatMessage](512, nil, MessageExpiration)
	a.llmReplyLru = expirable.NewLRU[string, struct{}](LlmReplyLimitCount, nil, LlmReplyLimitDuration)
	a.probabilityLlmTriggerRandom = rand.New(rand.NewSource(time.Now().UnixNano()))

	a.isLiving = true
	a.isLlmProcessing = false

	giftTimerMap := make(map[string]*GiftWithTimer)
	var giftTimerMapMutex sync.RWMutex

	// 消息处理 Handle
	dispatcherHandleMap := basic.DispatcherHandleMap{
		proto.OperationMessage: func(_ *basic.WsClient, msg *proto.Message) error {
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
					a.writeResultOK(ResultTypeDanmu, danmuData)

					go a.setUser(u)

					a.historyMsgLru.Add(d.MsgID, &ChatMessage{
						OpenId:    danmuData.OpenID,
						User:      danmuData.Uname,
						Message:   danmuData.Msg,
						Timestamp: time.Now(),
					})

					pitchRate := 0
					//if !livingCfg.DisableLlm {
					//	pitchRate = -100
					//}
					a.pushTTS(&tts.NewTaskParams{
						Text:      fmt.Sprintf("%s说：%s", d.Uname, d.Msg),
						PitchRate: pitchRate,
					}, false)

					if a.isLlmProcessing {
						break
					}

					if (danmuData.FansMedalWearingStatus &&
						danmuData.FansMedalName == FansMedalName &&
						danmuData.FansMedalLevel >= LlmReplyFansMedalLevel) || // 带10级粉丝牌
						danmuData.GuardLevel > 0 || // 舰长
						(danmuData.Uname == "巫女酱子" || danmuData.Uname == "青云-_-z") {
						a.startLlmReply(false)
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
					a.writeResultOK(ResultTypeSuperChat, scData)

					go a.setUser(u)

					a.historyMsgLru.Add(d.MsgID, &ChatMessage{
						OpenId:    scData.OpenID,
						User:      scData.Uname,
						Message:   scData.Msg,
						Timestamp: time.Now(),
					})
					a.pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("谢谢%s酱的醒目留言：%s", d.Uname, d.Message),
					}, false)
					a.startLlmReply(true)
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
					a.writeResultOK(ResultTypeGift, &GiftData{
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

					go a.setUser(u)

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
					go func(gt *GiftWithTimer) {
						defer gt.Timer.Stop()
						<-gt.Timer.C

						giftTimerMapMutex.Lock()
						delete(giftTimerMap, key)
						giftTimerMapMutex.Unlock()

						giftNum := atomic.LoadInt32(&gt.GiftNum)
						a.pushTTS(&tts.NewTaskParams{
							Text: fmt.Sprintf("谢谢%s酱赠送的%d个%s 么么哒", gt.Uname, giftNum, gt.GiftName),
						}, false)
					}(gt)
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
					a.writeResultOK(ResultTypeGuard, &GuardData{
						UserData:   u,
						GuardLevel: d.GuardLevel,
						GuardNum:   d.GuardNum,
						GuardUnit:  d.GuardUnit,
						Timestamp:  d.Timestamp,
						MsgID:      d.MsgID,
					})
					go a.setUser(u)
					guardName := getGuardLevelName(d.GuardLevel)
					a.pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("谢谢%s酱赠送的%d个%s%s，么么哒", d.UserInfo.Uname, d.GuardNum, d.GuardUnit, guardName),
					}, false)
					break
				}
			case *proto.CmdLiveStartData:
				{
					pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("主人开始直播啦，弹幕姬启动！直播分区为%s，直播间标题为%s", d.AreaName, d.Title),
					}, true)
					a.isLiving = true
					break
				}
			case *proto.CmdLiveEndData:
				{
					pushTTS(&tts.NewTaskParams{
						Text: "主人直播结束啦，今天辛苦了！",
					}, true)
					a.isLiving = false
					break
				}
			case *proto.CmdLiveRoomEnterData:
				{
					u := UserData{
						OpenID: d.OpenID,
						Uname:  d.Uname,
						UFace:  d.Uface,
					}
					a.writeResultOK(ResultTypeEnterRoom, &RoomEnterData{
						UserData:  u,
						Timestamp: d.Timestamp,
					})

					a.lastEnterUser = &u

					go func(openId string) {
						u, err := a.Dao.GetUser(context.Background(), openId)
						if err != nil {
							log.Errorf("GetUser open_id: %s err: %v", openId, err)
							return
						}

						if u == nil {
							return
						}

						if (u.FansMedalWearingStatus && u.FansMedalLevel >= RoomEnterTTSFansMedalLevel) ||
							u.GuardLevel > 0 {

							name := d.Uname
							if u.GuardLevel > 0 {
								guardName := getGuardLevelName(u.GuardLevel)
								name = guardName + name
							}

							pushTTS(&tts.NewTaskParams{
								Text: fmt.Sprintf("欢迎%s酱来到直播间", name),
							}, false)
						}
					}(d.OpenID)

					break
				}
			case *proto.CmdRoomChangeData:
				{
					pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("直播间信息发生变更，直播分区为%s，直播间标题为%s", d.AreaName, d.Title),
					}, true)
					break
				}
			case *proto.CmdInteractWordData:
				{
					a.writeResultOK(ResultTypeInteractWord, &InteractWordData{
						MsgID:     d.MsgID,
						OpenID:    d.OpenID,
						RoomID:    d.RoomID,
						Timestamp: d.Timestamp,
						Uname:     d.Uname,
					})
					pushTTS(&tts.NewTaskParams{
						Text: fmt.Sprintf("谢谢%s酱关注直播间", d.Uname),
					}, true)
					break
				}
			default:
				{
					break
				}
			}

			return nil
		},
	}

	a.wcs, err = basic.StartWebsocket(
		startResp,
		dispatcherHandleMap,
		onCloseHandle,
		a.slog,
	)
	if err != nil {
		log.Errorf("basic.StartWebsocket err: %v", err)
		a.writeResultError(ResultTypeRoom, CodeInternalError, err.Error())
		return
	}

	log.Infof("room_info: %v", startResp.AnchorInfo)
	a.writeResultOK(ResultTypeRoom, &RoomData{
		RoomID: startResp.AnchorInfo.RoomID,
		Uname:  startResp.AnchorInfo.Uname,
		UFace:  startResp.AnchorInfo.UFace,
	})
}

func (a *App) startLlmReply(force bool) {
	if !a.isLiving || a.livingCfg.DisableLlm {
		return
	}

	var msgs []*ChatMessage
	userMap := map[string]struct{}{}
	probabilityLlmTriggerCounter := -1 // 当前尝试触发的用户不算，所以初始值为-1
	for _, msg := range a.historyMsgLru.Values() {
		if time.Since(msg.Timestamp) <= LlmHistoryDuration {
			msgs = append(msgs, msg)
		}
		if time.Since(msg.Timestamp) <= DisableLlmByUserCountDuration {
			userMap[msg.OpenId] = struct{}{}
		}
		if time.Since(msg.Timestamp) <= ProbabilityLlmTriggerDuration {
			probabilityLlmTriggerCounter++
		}
	}

	if !force {
		llmReplyLruLen := a.llmReplyLru.Len()
		if llmReplyLruLen >= LlmReplyLimitCount {
			log.Infof("disable llm by reply count: %d", llmReplyLruLen)
			return
		}

		if len(userMap) >= DisableLlmByUserCount {
			log.Infof("disable llm by user count: %d", len(userMap))
			return
		}

		currentMsg := msgs[len(msgs)-1]
		if IsRepeatedChar(currentMsg.Message) {
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

		r := a.probabilityLlmTriggerRandom.Float64()
		fmt.Printf("r: %.2f, probability: %.2f\n", r, probability)
		if r <= probability {
			log.Infof("disable llm by probability: %.2f, counter: %d, compare: %.2f", r, probabilityLlmTriggerCounter, probability)
			return
		}
	}

	a.isLlmProcessing = true
	go func(msgs []*ChatMessage) {
		defer func() {
			a.isLlmProcessing = false
		}()

		llmMsgs := make([]*llm.ChatMessage, len(msgs))
		for i, msg := range msgs {
			llmMsgs[i] = &llm.ChatMessage{
				User:    msg.User,
				Message: msg.Message,
			}
		}
		llmRes, err := a.LLM.ChatWithLLM(context.Background(), llmMsgs)
		if err != nil {
			a.writeResultError(ResultTypeLLM, CodeInternalError, err.Error())
			log.Errorf("ChatWithLLM err: %v", err)
			return
		}
		a.writeResultOK(ResultTypeLLM, map[string]interface{}{
			"llm_result": llmRes,
		})
		a.llmReplyLru.Add(uuid.NewV4().String(), struct{}{})
		a.pushTTS(&tts.NewTaskParams{
			Text: llmRes,
		}, false)
	}(msgs)
}

func (a *App) StopConn() {
	if a.wcs != nil {
		a.wcs.Close()
		a.wcs = nil
	}
	a.startResp = nil
	a.connContext = nil
	if a.connCancel != nil {
		a.connCancel()
		a.connCancel = nil
	}
	if a.tk != nil {
		a.tk.Stop()
		a.tk = nil
	}
	a.lastEnterUser = nil
	if a.lastEnterUserTimer != nil {
		a.lastEnterUserTimer.Stop()
		a.lastEnterUserTimer = nil
	}
	if a.ttsQueue != nil {
		a.ttsQueue.Close()
		a.ttsQueue = nil
	}
}

func (a *App) setUser(userData UserData) {
	err := a.Dao.CreateOrUpdateUser(context.Background(), &dao.User{
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
