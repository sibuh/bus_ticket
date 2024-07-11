package ticket

import (
	"bytes"
	"context"
	"encoding/json"
	"event_ticket/internal/model"
	"event_ticket/internal/platform"
	"io/ioutil"
	"net/http"

	"golang.org/x/exp/slog"
)

type MockPaymentGateWay struct {
	logger *slog.Logger
	url    string
	platform.PaymentGatewayIntegrator
}

func InitMockGateway(logger *slog.Logger, url string) *MockPaymentGateWay {
	return &MockPaymentGateWay{logger: logger, url: url}
}

func (m *MockPaymentGateWay) CreateCheckoutSession(ctx context.Context, ticketInfo model.Ticket) (model.Session, error) {
	b, err := json.Marshal(ticketInfo)
	if err != nil {
		return model.Session{}, err
	}
	resp, err := http.Post(m.url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return model.Session{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.Session{}, err
	}
	r := model.Session{}
	if err := json.Unmarshal(body, &r); err != nil {
		return model.Session{}, err
	}
	return r, nil
}
