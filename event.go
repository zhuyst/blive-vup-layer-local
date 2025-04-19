package main

import (
	"blive-vup-layer/util"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type EventResult struct {
	Type string      `json:"type"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (s *Service) writeResultOK(resultType string, data interface{}) {
	s.writeResult(&EventResult{
		Type: resultType,
		Code: CodeOK,
		Msg:  "success",
		Data: data,
	})
}

func (s *Service) writeResultError(resultType string, code int, msg string) {
	s.writeResult(&EventResult{
		Type: resultType,
		Code: code,
		Msg:  msg,
	})
}

func (s *Service) writeResult(res *EventResult) {
	l := log.WithField("type", res.Type)
	msg, _ := json.Marshal(res)
	if res.Code == CodeOK {
		if res.Type != ResultTypeHeartbeat {
			l.Infof("write result type: %s, code: %d, data: %s", res.Type, res.Code, msg)
		}
	} else {
		errMsg := fmt.Sprintf("write result type: %s, code: %d, data: %s", res.Type, res.Code, msg)
		logLevel := log.WarnLevel
		if res.Code >= CodeInternalError {
			logLevel = log.ErrorLevel
			util.ShowErrorDialog(errMsg)
		}
		l.Log(logLevel, errMsg)
	}
	s.app.App.EmitEvent(res.Type, res)
}
