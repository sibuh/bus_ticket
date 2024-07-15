package session

import (
	"context"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"event_ticket/internal/storage"
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
		TicketID:      sess.ID,
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
