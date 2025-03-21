package main

import (
	"encoding/json"
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
	msg, _ := json.Marshal(res)
	if res.Code == CodeOK {
		if res.Type != ResultTypeHeartbeat {
			log.Infof("write result type: %s, code: %d, data: %s", res.Type, res.Code, msg)
		}
	} else {
		log.Errorf("write result type: %s, code: %d, data: %s", res.Type, res.Code, msg)
	}
	s.app.EmitEvent(res.Type, res)
}
