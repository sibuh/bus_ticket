package ticket

import (
	"context"
	"database/sql"
	"errors"
	"event_ticket/internal/constant"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"event_ticket/internal/module/schedule"
	"event_ticket/internal/platform"
	"event_ticket/internal/storage"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

type ticket struct {
	log           *slog.Logger
	storageTicket storage.Ticket
	platform      platform.PaymentGatewayIntegrator
	scheduler     *schedule.Scheduler
	db.Querier
}

type TicketStatus string

const (
	Reserved TicketStatus = "Reserved"
	Free     TicketStatus = "Free"
	Onhold   TicketStatus = "Onhold"
)

type PaymentStatus string

const (
	Succeeded PaymentStatus = "SUCCEEDED"
	Failed    PaymentStatus = "FAILED"
	Cancelled PaymentStatus = "CANCELLED"
	Pending   PaymentStatus = "PENDING"
)

func Init(log *slog.Logger, platform platform.PaymentGatewayIntegrator, q db.Querier, sc *schedule.Scheduler) module.Ticket {
	return &ticket{
		log:       log,
		platform:  platform,
		Querier:   q,
		scheduler: sc,
	}
}

func (t *ticket) ReserveTicket(ctx context.Context, req model.ReserveTicketRequest) (model.Session, error) {

	tkt, err := t.GetTicket(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newError := model.Error{
				ErrCode:   http.StatusNotFound,
				Message:   "the requested ticket is not found",
				RootError: err,
			}
			return model.Session{}, &newError
		}
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to get ticket",
			RootError: err,
		}
		return model.Session{}, &newError
	}
	if tkt.Status == string(Reserved) {
		newError := model.NewError(http.StatusBadRequest,
			"ticket is already reserved please try to reserve free ticket",
			fmt.Errorf("ticket reserved"))
		t.log.Error(newError.Error(), newError)

		return model.Session{}, newError
	}

	if tkt.Status == string(Onhold) {
		newError := model.NewError(http.StatusBadRequest,
			"ticket is onhold please try later",
			fmt.Errorf("ticket held"))
		t.log.Error(newError.Error(), newError)
		return model.Session{}, newError
	}

	tkt, err = t.UpdateTicketStatus(ctx, db.UpdateTicketStatusParams{
		ID:     req.ID,
		Status: req.Status,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newError := model.Error{
				ErrCode:   http.StatusNotFound,
				Message:   "ticket to unhold does not exist",
				RootError: err,
			}
			t.log.Error("ticket to unhold not found", newError)
			return model.Session{}, &newError
		}

		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to unhold ticket",
			RootError: err,
		}
		t.log.Error("failed to unhold ticket when checkout session creation fails", newError)
		return model.Session{}, &newError
	}
	if tkt.Status != string(Onhold) {
		newError := model.NewError(http.StatusInternalServerError, "ticket is not held successfully", nil)
		t.log.Error(newError.Error(), newError)
		return model.Session{}, newError
	}
	session, err := t.platform.CreateCheckoutSession(model.Ticket{
		ID:       tkt.ID,
		TripID:   tkt.TripID,
		TicketNo: tkt.TicketNo,
		BusNo:    tkt.BusNo,
		Status:   tkt.Status,
	})
	if err != nil {
		newError := model.NewError(http.StatusInternalServerError, "failed to create checkout session", err)
		t.log.Error(newError.Error(), newError)
		//unhold ticket if create checkout session fails
		_, err = t.UpdateTicketStatus(ctx, db.UpdateTicketStatusParams{
			ID:     tkt.ID,
			Status: string(constant.Free),
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				newError := model.Error{
					ErrCode:   http.StatusNotFound,
					Message:   "ticket to unhold does not exist",
					RootError: err,
				}
				t.log.Error("ticket to unhold not found", newError)
				return model.Session{}, &newError
			}

			newError := model.Error{
				ErrCode:   http.StatusInternalServerError,
				Message:   "failed to unhold ticket",
				RootError: err,
			}
			t.log.Error("failed to unhold ticket when checkout session creation fails", newError)
		}

		return model.Session{}, newError
	}
	storedSession, err := t.StoreCheckoutSession(ctx, db.StoreCheckoutSessionParams{
		ID:            session.ID,
		TicketID:      session.TicketID,
		PaymentStatus: session.PaymentStatus,
		PaymentURL:    session.PaymentURL,
		CancelURL:     session.CancelURL,
		Amount:        session.Amount,
		CreatedAt:     session.CreatedAt,
	})
	if err != nil {
		newError := model.NewError(http.StatusInternalServerError, "failed to store checkout session", err)
		t.log.Error(newError.Error(), newError)
		return model.Session{}, newError
	}

	sId := storedSession.ID
	ch := make(chan string)

	go t.scheduler.Schedule(sId, ch, 10*time.Minute, t.QueryFunc)
	return model.Session{
		ID:            storedSession.ID,
		TicketID:      storedSession.TicketID,
		PaymentStatus: storedSession.PaymentStatus,
		PaymentURL:    storedSession.PaymentURL,
		CancelURL:     storedSession.CancelURl,
		Amount:        storedSession.Amount,
		CreatedAt:     storedSession.CreatedAt,
	}, err
}
func (t *ticket) QueryFunc(id string) error {
	for i := 0; i < 5; i++ {

		d := time.Duration(i)
		sleep := 30 * time.Second * (1 + d)

		status, err := t.platform.CheckPaymentStatus(context.Background(), id)
		if err != nil {
			t.log.Error("payment status request failed", err)
			time.Sleep(sleep)
			continue
		}
		if status == string(Succeeded) || status == string(Failed) {
			_, err := t.storageTicket.UpdateTicket(context.Background(), model.ReserveTicketRequest{ID: id, Status: status})
			if err != nil {
				t.log.Error("failed to update ticket status after cancelling session")
			}
			break
		} else if status == string(Pending) {
			cancel, err := t.platform.CancelCheckoutSession(context.Background(), id)
			if err != nil {
				t.log.Error("failed to cancel check out session", err)
				time.Sleep(sleep)
				continue
			}
			if cancel {
				_, err := t.storageTicket.UpdateTicket(context.Background(), model.ReserveTicketRequest{ID: id, Status: string(Free)})
				if err != nil {
					t.log.Error("failed to update ticket status after cancelling session")
				}
				break
			}
		}

	}

	return nil
}
