package ticket

import (
	"context"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"event_ticket/internal/platform"
	"event_ticket/internal/storage"
	"net/http"

	"golang.org/x/exp/slog"
)

type ticket struct {
	log           *slog.Logger
	storageTicket storage.Ticket
	platform      platform.PaymentGatewayIntegrator
	session       storage.Session
}

type TicketStatus string

const (
	Reserved TicketStatus = "Reserved"
	Free     TicketStatus = "Free"
	Onhold   TicketStatus = "Onhold"
)

func Init(log *slog.Logger, tkt storage.Ticket, platform platform.PaymentGatewayIntegrator, ssn storage.Session) module.Ticket {
	return &ticket{
		log:           log,
		storageTicket: tkt,
		platform:      platform,
		session:       ssn,
	}
}

func (t *ticket) ReserveTicket(ctx context.Context, req model.ReserveTicketRequest, scheduler func()) (model.Session, error) {
	tkt, err := t.storageTicket.GetTicket(ctx, req.ID)
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

	tkt, err = t.storageTicket.HoldTicket(ctx, req)

	if err != nil {
		return model.Session{}, err
	}
	if tkt.Status != string(Onhold) {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "ticket is not held successfully",
			RootError: nil,
		}
		t.log.Error(newError.Error(), newError)
		return model.Session{}, &newError
	}
	session, err := t.platform.CreateCheckoutSession(tkt)
	if err != nil {
		//unhold ticket if create checkout session fails
		_, err = t.storageTicket.UnholdTicket(tkt.ID)
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
	storedSession, err := t.session.StoreCheckoutSession(ctx, session)
	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to store checkout session",
			RootError: err,
		}
		t.log.Error(newError.Error(), newError)
		return model.Session{}, &newError
	}
	scheduler()
	return storedSession, err
}
