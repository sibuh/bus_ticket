package platform

import (
	"context"
	"event_ticket/internal/model"
)

type PaymentGatewayIntegrator interface {
	CreateCheckoutSession(ticketInfo model.Ticket) (model.Session, error)
	CancelCheckoutSession(ctx context.Context, ID string) (bool, error)
}
