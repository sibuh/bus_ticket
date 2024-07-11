package platform

import (
	"context"
	"event_ticket/internal/model"
)

type PaymentGatewayIntegrator interface {
	CreateCheckoutSession(ticketInfo model.Ticket)
	CancelCheckoutSession(ctx context.Context, ID string) (bool, error)
}
