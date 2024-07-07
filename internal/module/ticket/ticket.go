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
	log           *slog.Logger
	storageTicket storage.Ticket
	platform      platform.PaymentGatewayIntegrator
}

type TicketStatus string

const (
	Reserved TicketStatus = "Reserved"
	Free     TicketStatus = "Free"
	Onhold   TicketStatus = "Onhold"
)

func Init(log *slog.Logger, tkt storage.Ticket, platform platform.PaymentGatewayIntegrator) module.Ticket {
	return &ticket{
		log:           log,
		storageTicket: tkt,
		platform:      platform,
	}
}

func (t *ticket) ReserveTicket(ctx context.Context, tktNo, tripId int32) (model.Session, error) {
	tkt, err := t.storageTicket.GetTicket(tktNo, tripId)
	if err != nil {
		return model.Session{}, err
	}
	if tkt.Status == string(Reserved) {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "ticket is already reserved please try to reserve free ticket",
			RootError: nil,
		}
		return model.Session{}, &newError
	}

	if tkt.Status == string(Onhold) {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "ticket is onhold please try later",
			RootError: nil,
		}
		return model.Session{}, &newError
	}

	tkt, err = t.storageTicket.HoldTicket(tktNo, tripId)

	if err != nil {
		return model.Session{}, err
	}
	if tkt.Status != string(Onhold) {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "ticket is not held successfully",
			RootError: nil,
		}
		return model.Session{}, &newError
	}
	checkoutUrl, err := t.platform.CreateCheckoutSession(ctx, tkt)
	if err != nil {
		//unhold ticket if create checkout session fails
		_, err = t.storageTicket.UnholdTicket(tktNo, tripId)
		if err != nil {
			newError := model.Error{
				ErrCode:   http.StatusInternalServerError,
				Message:   "failed to unhold ticket",
				RootError: err,
			}
			t.log.Error("failed to unhold ticket when creating checkout session fails", newError)
		}

		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to create checkout session",
			RootError: err,
		}

		t.log.Error("failed to create checkout session", newError)
		return model.Session{}, &newError
	}

	// SERVIER.ON('INIT, HANDLESERVERINIT)
	// READ FROM DATABASE PENDING STATUS CHECKOUT SESSION `[]SESSION`
	// LOOP TIME.AFTERfUNC(TIME.NOW() - SESSION.TIME)
	time.AfterFunc(time.Second, func() {
		func(tktNo, tripId int32, logger *slog.Logger) {
			_, err := t.platform.CancelCheckoutSession(ctx, "")
			if err != nil {
				newError := model.Error{
					ErrCode:   http.StatusInternalServerError,
					Message:   "failed to cancel checkout session",
					RootError: err,
				}
				t.log.Error(newError.Error(), newError.ErrCode)
			}
			_, err = t.storageTicket.UnholdTicket(tktNo, tripId)
			if err != nil {
				logger.Error("failed to unhold ticket", err)
			}
		}(tktNo, tripId, t.log)
	},
	)
	return checkoutUrl, nil
}
