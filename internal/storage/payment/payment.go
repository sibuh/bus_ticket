package payment

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

type payment struct {
	logger *slog.Logger
	db     *db.Queries
}

func Init(logger *slog.Logger, db *db.Queries) storage.Payment {
	return &payment{
		logger: logger,
		db:     db,
	}
}

func (p *payment) RecordPaymentIntent(ctx context.Context, param model.CreateIntentParam) (model.Payment, error) {
	payment, err := p.db.RecordPayment(ctx, db.RecordPaymentParams{
		UserID:        param.UserID,
		EventID:       param.EventID,
		IntentID:      param.IntentID,
		PaymentStatus: "pending",
		CheckInStatus: "pending",
	})
	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to record paynent intent",
			RootError: err,
		}
		p.logger.Error("failed to record payment intent", err)
		return model.Payment{}, &newError
	}
	return model.Payment(payment), nil
}

func (p *payment) GetPayment(ctx context.Context, intentID string) (model.Payment, error) {
	pmt, err := p.db.GetPayment(ctx, intentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newError := model.Error{
				ErrCode:   http.StatusNotFound,
				Message:   fmt.Sprintf("no payment data found for IntentID %s", intentID),
				RootError: err,
			}
			p.logger.Info(fmt.Sprintf("no payment data found for IntentID %s", intentID))
			return model.Payment{}, &newError
		}
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to get payment information",
			RootError: err,
		}
		p.logger.Error("failed to get payment information")
		return model.Payment{}, &newError

	}
	return model.Payment(pmt), nil
}
