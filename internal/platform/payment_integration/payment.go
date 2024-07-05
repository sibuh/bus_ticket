package paymentintegration

import (
	"context"
	"event_ticket/internal/model"
	"event_ticket/internal/platform"

	"golang.org/x/exp/slog"
)

type paymentGateway struct {
	loger *slog.Logger
}

func Init(logger *slog.Logger) platform.PaymentGatewayIntegrator {
	return &paymentGateway{}
}
func (p *paymentGateway) CreateCheckoutSession(ctx context.Context, ticketInfo model.Ticket) (model.Session, error) {
	//TODO:create checkout session
	return model.Session{}, nil
}
func (p *paymentGateway) CancelCheckoutSession(ctx context.Context, sessionId string) (bool, error) {
	return true, nil
}
