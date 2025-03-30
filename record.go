package main

import (
	"context"
	"encoding/base64"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

func (s *Service) CommitRecord(base64Str string) {
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		log.Errorf("base64 decoding error: %s", err.Error())
		s.writeResultError(ResultTypeRecordResult, CodeBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := s.Sr.RecognitionWav(ctx, decodedBytes)
	if err != nil {
		log.Errorf("RecognitionWav error: %s", err.Error())
		s.writeResultError(ResultTypeRecordResult, CodeBadRequest, err.Error())
		return
	}
	log.Infof("RecognitionWav result: %s", res)
	s.writeResultOK(ResultTypeRecordResult, map[string]interface{}{
		"result": res,
	})

	userData := UserData{
		Uname: "巫女酱子",
		UFace: "https://i0.hdslb.com/bfs/face/c4ebadbf1926960266223a3a0508989372b67ecf.jpg",
	}
	msgId := uuid.NewV4().String()
	s.writeResultOK(ResultTypeDanmu, &DanmuData{
		UserData: userData,
		Msg:      res,
		MsgID:    msgId,
	})
	s.historyMsgLru.Add(msgId, &ChatMessage{
		User:      &userData,
		Message:   res,
		Timestamp: time.Now(),
	})
	s.startLlmReply(true)
}
