package ticket

import (
	"context"
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
	session       storage.Session
	scheduler     *schedule.Scheduler
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

func Init(log *slog.Logger, tkt storage.Ticket, platform platform.PaymentGatewayIntegrator, ssn storage.Session, sc *schedule.Scheduler) module.Ticket {
	return &ticket{
		log:           log,
		storageTicket: tkt,
		platform:      platform,
		session:       ssn,
		scheduler:     sc,
	}
}

func (t *ticket) ReserveTicket(ctx context.Context, req model.ReserveTicketRequest) (model.Session, error) {
	tkt, err := t.storageTicket.GetTicket(ctx, req.ID)
	if err != nil {
		return model.Session{}, err
	}
	if tkt.Status == string(Reserved) {
		newError := model.NewError(http.StatusBadRequest, "ticket is already reserved please try to reserve free ticket", fmt.Errorf("ticket reserved"))
		t.log.Error(newError.Error(), newError)

		return model.Session{}, newError
	}

	if tkt.Status == string(Onhold) {
		newError := model.NewError(http.StatusBadRequest, "ticket is onhold please try later", fmt.Errorf("ticket held"))
		t.log.Error(newError.Error(), newError)
		return model.Session{}, newError
	}

	tkt, err = t.storageTicket.UpdateTicket(ctx, req)

	if err != nil {
		return model.Session{}, err
	}
	if tkt.Status != string(Onhold) {
		newError := model.NewError(http.StatusInternalServerError, "ticket is not held successfully", nil)
		t.log.Error(newError.Error(), newError)
		return model.Session{}, newError
	}
	session, err := t.platform.CreateCheckoutSession(tkt)
	if err != nil {
		newError := model.NewError(http.StatusInternalServerError, "failed to create checkout session", err)
		t.log.Error(newError.Error(), newError)
		//unhold ticket if create checkout session fails
		_, err = t.storageTicket.UnholdTicket(tkt.ID)
		if err != nil {
			newError := model.NewError(http.StatusInternalServerError, "failed to unhold ticket", err)
			t.log.Error(newError.Error(), newError)
		}

		return model.Session{}, newError
	}
	storedSession, err := t.session.StoreCheckoutSession(ctx, session)
	if err != nil {
		newError := model.NewError(http.StatusInternalServerError, "failed to store checkout session", err)
		t.log.Error(newError.Error(), newError)
		return model.Session{}, newError
	}

	sId := storedSession.ID
	ch := make(chan string)

	go t.scheduler.Schedule(sId, ch, 10*time.Minute, t.QueryFunc)
	return storedSession, err
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
