package payment

import (
	"context"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"event_ticket/internal/storage"
	"net/http"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"golang.org/x/exp/slog"
)

type payment struct {
	logger *slog.Logger
	es     storage.Event
}

func Init(logger *slog.Logger, es storage.Event) module.Payment {
	return &payment{
		logger: logger,
		es:     es,
	}
}

func (p *payment) CreatePaymentIntent(ctx context.Context, userID, eventID int32) (string, error) {
	event, err := p.es.FetchEvent(ctx, eventID)
	if err != nil {
		return "", err
	}
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(2000 + event.Price)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to create payment intent",
			RootError: err,
		}
		p.logger.Error("failed to create stripe payment intent")
		return "", &newError
	}
	//TODO register the user session into payment database
	return pi.ClientSecret, nil
}
