package ticket

import (
	"context"
	"database/sql"
	"errors"
	"event_ticket/internal/constant"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"event_ticket/internal/storage"
	"net/http"

	"golang.org/x/exp/slog"
)

type ticket struct {
	logger *slog.Logger
	db     db.Querier
}

func Init(logger *slog.Logger, db db.Querier) storage.Ticket {
	return &ticket{
		logger: logger,
		db:     db,
	}
}
func (t *ticket) HoldTicket(ctx context.Context, req model.ReserveTicketRequest) (model.Ticket, error) {
	tkt, err := t.db.UpdateTicketStatus(ctx, db.UpdateTicketStatusParams{
		ID:     req.ID,
		Status: string(constant.Onhold),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newError := model.Error{
				ErrCode:   http.StatusNotFound,
				Message:   "the requested ticket is not found",
				RootError: err,
			}
			return model.Ticket{}, &newError
		}
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to update ticket status",
			RootError: err,
		}
		return model.Ticket{}, &newError
	}
	return model.Ticket{
		TripID:   tkt.TripID,
		BusNo:    tkt.BusNo,
		TicketNo: tkt.TicketNo,
		Status:   tkt.Status,
	}, nil
}
func (t *ticket) GetTicket(ctx context.Context, id string) (model.Ticket, error) {
	tkt, err := t.db.GetTicket(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newError := model.Error{
				ErrCode:   http.StatusNotFound,
				Message:   "the requested ticket is not found",
				RootError: err,
			}
			return model.Ticket{}, &newError
		}
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to get ticket",
			RootError: err,
		}
		return model.Ticket{}, &newError
	}

	return model.Ticket{
		ID:       tkt.ID,
		TripID:   tkt.TripID,
		TicketNo: tkt.TicketNo,
		BusNo:    tkt.BusNo,
		Status:   tkt.Status,
	}, nil
}
func (t *ticket) UnholdTicket(ID string) (model.Ticket, error) {
	tkt, err := t.db.UpdateTicketStatus(context.Background(), db.UpdateTicketStatusParams{
		ID:     ID,
		Status: string(constant.Free),
	})
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			newError := model.Error{
				ErrCode:   http.StatusNotFound,
				Message:   "ticket to unhold does not exist",
				RootError: err,
			}
			t.logger.Error("ticket to unhold not found", newError)
			return model.Ticket{}, &newError
		}

		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to unhold ticket",
			RootError: err,
		}
		t.logger.Error("failed to unhold ticket when checkout session creation fails", newError)
		return model.Ticket{}, &newError
	}
	return model.Ticket{
		ID:       tkt.ID,
		TripID:   tkt.TripID,
		TicketNo: tkt.TicketNo,
		BusNo:    tkt.BusNo,
		Status:   tkt.Status,
	}, nil
}
