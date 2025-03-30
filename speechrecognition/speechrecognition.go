package speechrecognition

import (
	"blive-vup-layer/config"
	"bytes"
	"errors"
	"github.com/aliyun/alibabacloud-nls-go-sdk"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"io"
	syslog "log"
	"time"
)

type Task struct {
	TaskId string
	Logger *log.Entry
	Result string
	Err    error

	sr       *nls.SpeechRecognition
	param    nls.SpeechRecognitionStartParam
	pcmAudio []byte
}

func (task *Task) onTaskFailed(text string, param interface{}) {
	task.Logger.Println("TaskFailed:", text)
	task.Result = text
}

func (task *Task) onStarted(text string, param interface{}) {
	task.Logger.Println("onStarted:", text)
	task.Result = text
}

func (task *Task) onResultChanged(text string, param interface{}) {
	task.Logger.Println("onResultChanged:", text)
	task.Result = text
}

func (task *Task) onCompleted(text string, param interface{}) {
	task.Logger.Println("onCompleted:", text)
	task.Result = text
}

func (task *Task) onClose(param interface{}) {
	task.Logger.Println("onClosed:")
}

func (task *Task) waitReady(ch chan bool) error {
	select {
	case done := <-ch:
		{
			if !done {
				task.Logger.Println("Wait failed")
				return errors.New("wait failed")
			}
			task.Logger.Println("Wait done")
		}
	case <-time.After(20 * time.Second):
		{
			task.Logger.Println("Wait timeout")
			return errors.New("wait timeout")
		}
	}
	return nil
}

type SpeechRecognition struct {
	cfg *config.AliyunTTSConfig
}

func NewSpeechRecognition(cfg *config.AliyunTTSConfig) *SpeechRecognition {
	return &SpeechRecognition{
		cfg: cfg,
	}
}

type NewTaskParams struct {
	PcmAudio []byte
}

func (sr *SpeechRecognition) NewTask(params *NewTaskParams) (*Task, error) {
	nlsCfg, err := nls.NewConnectionConfigWithAKInfoDefault(
		nls.DEFAULT_URL,
		sr.cfg.AppKey, sr.cfg.AccessKey, sr.cfg.SecretKey,
	)
	if err != nil {
		log.Errorf("NewConnectionConfigWithAKInfoDefault err: %v", err)
		return nil, err
	}

	taskId := uuid.NewV4().String()
	l := log.WithField("task_id", taskId)
	param := nls.SpeechRecognitionStartParam{
		Format:     "wav",
		SampleRate: 16000,
		//EnableIntermediateResult:       true,
		EnablePunctuationPrediction:    true,
		EnableInverseTextNormalization: true,
	}

	t := &Task{
		TaskId:   taskId,
		Logger:   l,
		param:    param,
		pcmAudio: params.PcmAudio,
	}

	nlsLog := nls.NewNlsLogger(io.Discard, "NLS", syslog.LstdFlags|syslog.Lmicroseconds)
	nlsSr, err := nls.NewSpeechRecognition(nlsCfg, nlsLog,
		t.onTaskFailed, t.onStarted, t.onResultChanged,
		t.onCompleted, t.onClose, nil)
	if err != nil {
		l.Errorf("NewSpeechRecognition err: %v", err)
		return nil, err
	}
	t.sr = nlsSr

	return t, nil
}

func (task *Task) Run() (string, error) {
	task.Logger.Println("Sr start")
	ready, err := task.sr.Start(task.param, nil)
	if err != nil {
		task.Logger.Errorf("Start err: %v", err)
		task.sr.Shutdown()
		task.Err = err
		return "", err
	}
	if err = task.waitReady(ready); err != nil {
		task.Logger.Errorf("waitReady err: %v", err)
		task.sr.Shutdown()
		task.Err = err
		return "", err
	}
	buffers := nls.LoadPcmInChunk(bytes.NewBuffer(task.pcmAudio), 320)
	for _, data := range buffers.Data {
		if data == nil {
			continue
		}
		if err := task.sr.SendAudioData(data.Data); err != nil {
			task.Logger.Errorf("SendAudioData err: %v", err)
			return "", err
		}
		time.Sleep(10 * time.Millisecond)
	}
	ready, err = task.sr.Stop()
	if err != nil {
		task.Logger.Errorf("Stop err: %v", err)
		task.sr.Shutdown()
		task.Err = err
		return "", err
	}

	if err = task.waitReady(ready); err != nil {
		task.Logger.Errorf("waitReady err: %v", err)
		task.sr.Shutdown()
		task.Err = err
		return "", err
	}

	task.Logger.Infof("Sr done")
	task.sr.Shutdown()

	return task.Result, nil
}
