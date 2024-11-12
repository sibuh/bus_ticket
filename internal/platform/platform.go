package platform

import (
	"bus_ticket/internal/model"
	"context"
)

type PaymentGatewayIntegrator interface {
	CreateCheckoutSession(ticketInfo model.Ticket) (model.Session, error)
	CancelCheckoutSession(ctx context.Context, ID string) (bool, error)
	CheckPaymentStatus(ctx context.Context, ID string) (string, error)
}
