package ticket

import (
	"context"
	"event_ticket/internal/module"
	"event_ticket/internal/platform"
	"event_ticket/internal/storage"

	"golang.org/x/exp/slog"
)

type ticket struct {
	log           slog.Logger
	storageTicket storage.Ticket
	platform      platform.PaymentGatewayIntegrator
}

func Init(log slog.Logger, tkt storage.Ticket, platform platform.PaymentGatewayIntegrator) module.Ticket {
	return &ticket{
		log:           log,
		storageTicket: tkt,
		platform:      platform,
	}
}
func (t *ticket) HoldTicket(ctx context.Context, tktNo, tripId int32) (string, error) {
	tkt, err := t.storageTicket.HoldTicket(tktNo, tripId)
	if err != nil {
		return "", err
	}
	url, err := t.platform.CreateCheckoutSession(ctx, tkt)
	if err != nil {
		return "", err
	}
	return url, nil
}
