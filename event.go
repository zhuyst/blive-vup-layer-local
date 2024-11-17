package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type EventResult struct {
	Type string      `json:"type"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (a *App) writeResultOK(resultType string, data interface{}) {
	a.writeResult(&EventResult{
		Type: resultType,
		Code: CodeOK,
		Msg:  "success",
		Data: data,
	})
}

func (a *App) writeResultError(resultType string, code int, msg string) {
	a.writeResult(&EventResult{
		Type: resultType,
		Code: code,
		Msg:  msg,
	})
}

func (a *App) writeResult(res *EventResult) {
	msg, _ := json.Marshal(res)
	if res.Code == CodeOK {
		if res.Type != ResultTypeHeartbeat {
			log.Infof("write result type: %s, code: %d, data: %s", res.Type, res.Code, msg)
		}
	} else {
		log.Errorf("write result type: %s, code: %d, data: %s", res.Type, res.Code, msg)
	}
	runtime.EventsEmit(a.appContext, res.Type, res)
}
