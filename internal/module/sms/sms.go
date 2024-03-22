package sms

import (
	"bytes"
	"encoding/json"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/exp/slog"
)

type sms struct {
	logger  slog.Logger
	token   string
	api     string
	message string
}

func Init(logger slog.Logger, token, api, msg string) module.Sms {
	return &sms{
		logger:  logger,
		token:   token,
		api:     api,
		message: msg,
	}
}

func (s *sms) SendSms(user db.User, wg *sync.WaitGroup) error {
	defer wg.Done()
	var payload = model.Sms{
		Token: s.token,
		Phone: user.Phone,
		Msg:   fmt.Sprintf(s.message, user.SessionID),
	}

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error("failed to marshal sms request body", err)
		return err
	}

	_, err = http.Post(s.api, "application/json", bytes.NewBuffer(payloadByte))
	if err != nil {
		s.logger.Error("post request to sms gate way failed", err)
		return err
	}
	return nil
}
