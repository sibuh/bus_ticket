package paymentintegration

import (
	"context"
	"event_ticket/internal/model"

	"golang.org/x/exp/slog"
)

type MockPaymentGateWay struct {
	logger slog.Logger
}

func InitMock(logger slog.Logger) *MockPaymentGateWay {
	return &MockPaymentGateWay{logger: logger}
}

func (m *MockPaymentGateWay) CreateCheckoutSession(ctx context.Context, ticketInfo model.Ticket) (url string, err error) {

	return "https://chapa.com/checkout/session-id", nil
}
