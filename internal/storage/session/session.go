package session

import (
	"context"
	"database/sql"
	"errors"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"event_ticket/internal/storage"
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"
)

type session struct {
	logger *slog.Logger
	db     db.Querier
}

func Init(logger *slog.Logger, db db.Querier) storage.Session {
	return &session{
		logger: logger,
		db:     db,
	}
}
func (s *session) StoreCheckoutSession(ctx context.Context, sess model.Session) (model.Session, error) {

	ssn, err := s.db.StoreCheckoutSession(ctx, db.StoreCheckoutSessionParams{
		ID:            sess.ID,
		TicketID:      sess.TicketID,
		PaymentStatus: sess.PaymentStatus,
		PaymentURL:    sess.PaymentURL,
		CancelURL:     sess.CancelURL,
		Amount:        sess.Amount,
		CreatedAt:     sess.CreatedAt,
	})

	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to store checkout session",
			RootError: err,
		}
		s.logger.Error("failed to store checkout session", newError)
		return model.Session{}, &newError
	}
	return model.Session{
		ID:            ssn.ID,
		TicketID:      ssn.TicketID,
		PaymentStatus: ssn.PaymentStatus,
		PaymentURL:    ssn.PaymentURL,
		CancelURL:     ssn.CancelURl,
		Amount:        ssn.Amount,
		CreatedAt:     ssn.CreatedAt,
	}, nil
}

func (s *session) GetTicketStatus(ctx context.Context, sid string) (string, error) {
	status, err := s.db.GetTicketStatus(ctx, sid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newError := model.Error{
				ErrCode:   http.StatusNotFound,
				Message:   fmt.Sprintf("session with %s not found", sid),
				RootError: err,
			}
			s.logger.Info(newError.Message, newError)
			return "", &newError
		}
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to get requested session",
			RootError: err,
		}
		s.logger.Error(newError.Error(), newError)
		return "", &newError
	}
	return status, nil
}
