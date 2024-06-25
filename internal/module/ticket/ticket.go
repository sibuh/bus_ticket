package ticket

import (
	"context"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"event_ticket/internal/platform"
	"event_ticket/internal/storage"
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

type ticket struct {
	log           slog.Logger
	storageTicket storage.Ticket
	platform      platform.PaymentGatewayIntegrator
}
type TicketStatus string

const (
	Reserved TicketStatus = "Reserved"
	Free     TicketStatus = "Free"
	Onhold   TicketStatus = "Onhold"
)

func Init(log slog.Logger, tkt storage.Ticket, platform platform.PaymentGatewayIntegrator) module.Ticket {
	return &ticket{
		log:           log,
		storageTicket: tkt,
		platform:      platform,
	}
}

func (t *ticket) ReserveTicket(ctx context.Context, tktNo, tripId int32) (string, error) {
	tkt, err := t.storageTicket.GetTicket(tktNo, tripId)
	if err != nil {
		return "", err
	}
	if tkt.Status == string(Reserved) {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "ticket is already reserved please try to reserve free ticket",
			RootError: nil,
		}
		return "", &newError
	}

	if tkt.Status == string(Onhold) {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "ticket is onhold please try later",
			RootError: nil,
		}
		return "", &newError
	}

	tkt, err = t.storageTicket.HoldTicket(tktNo, tripId)

	if err != nil {
		return "", err
	}
	if tkt.Status != string(Onhold) {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "ticket is not held successfully",
			RootError: nil,
		}
		return "", &newError
	}
	checkoutUrl, err := t.platform.CreateCheckoutSession(ctx, tkt)
	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to create checkout session",
			RootError: err,
		}
		return "", &newError
	}
	time.AfterFunc(time.Second, func() {
		func(tktNo, tripId int32, logger slog.Logger) {
			_, err := t.storageTicket.UnholdTicket(tktNo, tripId)
			if err != nil {
				logger.Error("failed to unhold ticket", err)
			}
		}(tktNo, tripId, t.log)
	},
	)
	return checkoutUrl, nil
}
