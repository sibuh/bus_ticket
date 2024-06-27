package platform

import (
	"context"
	"event_ticket/internal/model"
)

type PaymentGatewayIntegrator interface {
	CreateCheckoutSession(ctx context.Context, ticketInfo model.Ticket) (url string, err error)
	CancelCheckoutSession(ctx context.Context, sessionId string) (bool, error)
}
