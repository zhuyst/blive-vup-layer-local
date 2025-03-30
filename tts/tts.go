package tts

import (
	"blive-vup-layer/config"
	nls "blive-vup-layer/tts/alibabacloud-nls-go-sdk"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	syslog "log"
	"os"
	"path"
	"time"
)

type TTS struct {
	cfg *config.AliyunTTSConfig
}

func NewTTS(cfg *config.AliyunTTSConfig) (*TTS, error) {
	if err := os.MkdirAll(config.ResultFilePath, os.ModePerm); err != nil {
		return nil, err
	}
	return &TTS{cfg: cfg}, nil
}

type Task struct {
	TaskId string
	Logger *log.Entry

	Cache bool
	Fname string

	TmpFile  *os.File
	TmpFname string
	Err      error

	param           nls.SpeechSynthesisStartParam
	text            string
	speechSynthesis *nls.SpeechSynthesis
}

type NewTaskParams struct {
	Text      string
	PitchRate int
}

type TTSParams struct {
	nls.SpeechSynthesisStartParam
	Text string `json:"text"`
}

func (tts *TTS) NewTask(params *NewTaskParams) (*Task, error) {
	param := nls.SpeechSynthesisStartParam{
		Voice:      "voice-3e06127",
		Format:     "wav",
		SampleRate: 48000,
		Volume:     50,
		SpeechRate: -100,
		PitchRate:  params.PitchRate,
	}
	p := TTSParams{
		SpeechSynthesisStartParam: param,
		Text:                      params.Text,
	}
	j, _ := json.Marshal(p)
	md5Str := md5String(string(j))

	taskId := fmt.Sprintf("%s-%d", md5Str, len(params.Text))
	l := log.WithField("task_id", taskId)

	fname := path.Join(config.ResultFilePath, fmt.Sprintf("tts-%s.wav", taskId))
	if isFileExist(fname) {
		return &Task{
			TaskId: taskId,
			Logger: l,
			Fname:  fname,

			Cache: true,

			param: param,
			text:  params.Text,
		}, nil
	}

	tmpFname := path.Join(config.ResultFilePath, fmt.Sprintf("tts-%s.wav.tmp", taskId))
	os.Remove(tmpFname)
	tmpFout, err := os.OpenFile(tmpFname, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	t := &Task{
		TaskId:   taskId,
		Logger:   l,
		TmpFile:  tmpFout,
		TmpFname: tmpFname,
		Fname:    fname,

		Cache: false,

		param: param,
		text:  params.Text,
	}

	l.Infof("new tts: %s", t.text)
	nlsLog := nls.NewNlsLogger(io.Discard, "NLS", syslog.LstdFlags|syslog.Lmicroseconds)
	//nlsLog.SetDebug(true)

	nlsCfg, err := nls.NewConnectionConfigWithAKInfoDefault(
		nls.DEFAULT_URL,
		tts.cfg.AppKey, tts.cfg.AccessKey, tts.cfg.SecretKey,
	)
	if err != nil {
		l.Errorf("NewConnectionConfigWithAKInfoDefault err: %v", err)
		return nil, err
	}

	ss, err := nls.NewSpeechSynthesis(nlsCfg, nlsLog, false,
		t.onTaskFailed, t.onSynthesisResult, nil,
		t.onCompleted, t.onClose, param)
	if err != nil {
		l.Errorf("NewTask err: %v", err)
		return nil, err
	}

	t.speechSynthesis = ss
	return t, nil
}

func (task *Task) Run() (string, error) {
	if task.Cache {
		return task.Fname, nil
	}

	defer task.TmpFile.Close()

	ch, err := task.speechSynthesis.Start(task.text, task.param, nil)
	if err != nil {
		task.Logger.Errorf("Start err: %v", err)
		task.speechSynthesis.Shutdown()
		task.Err = err
		return "", err
	}

	err = task.waitReady(ch)
	if err != nil {
		task.Logger.Errorf("waitReady err: %v", err)
		task.speechSynthesis.Shutdown()
		task.Err = err
		return "", err
	}
	task.Logger.Infof("Synthesis done")
	task.speechSynthesis.Shutdown()

	os.Rename(task.TmpFname, task.Fname)

	return task.Fname, nil
}

func (task *Task) onTaskFailed(text string, param interface{}) {
	task.Logger.Errorf("TaskFailed: %s", text)
	task.TmpFile.Close()
}

func (task *Task) onSynthesisResult(data []byte, param interface{}) {
	task.TmpFile.Write(data)
}

func (task *Task) onCompleted(text string, param interface{}) {
	task.Logger.Infof("onCompleted: %s", text)
	task.TmpFile.Close()
}

func (task *Task) onClose(param interface{}) {
	task.Logger.Infof("onClosed")
	task.TmpFile.Close()
}

func (task *Task) waitReady(ch chan bool) error {
	select {
	case done := <-ch:
		{
			if !done {
				task.Logger.Error("wait failed")
				return errors.New("wait failed")
			}
			task.Logger.Debugf("Wait done")
		}
	case <-time.After(60 * time.Second):
		{
			task.Logger.Error("Wait timeout")
			return errors.New("wait timeout")
		}
	}
	return nil
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func md5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
