package paymentintegration

import (
	"context"
	"event_ticket/internal/model"
	"event_ticket/internal/platform"

	"golang.org/x/exp/slog"
)

type paymentGateway struct {
	logger *slog.Logger
	url    string
}

func Init(logger *slog.Logger, url string) platform.PaymentGatewayIntegrator {
	return &paymentGateway{
		logger: logger,
		url:    url,
	}
}
func (p *paymentGateway) CreateCheckoutSession(ctx context.Context, ticketInfo model.Ticket) (model.Session, error) {
	//TODO:create checkout session
	return model.Session{}, nil
}
func (p *paymentGateway) CancelCheckoutSession(ctx context.Context, sessionId string) (bool, error) {
	return true, nil
}
