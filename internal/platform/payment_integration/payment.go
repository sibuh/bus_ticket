package paymentintegration

import (
	"bytes"
	"context"
	"encoding/json"
	"event_ticket/internal/model"
	"event_ticket/internal/platform"
	"net/http"

	"golang.org/x/exp/slog"
)

type paymentGateway struct {
	logger *slog.Logger
	url    string
}

func Init(logger *slog.Logger, url string) platform.PaymentGatewayIntegrator {
	return &paymentGateway{
		logger: logger,
		url:    url,
	}
}
func (p *paymentGateway) CreateCheckoutSession(ticketInfo model.Ticket) {
	b, _ := json.Marshal(ticketInfo)
	// if err != nil {
	// 	return model.Session{}, err
	// }
	_, _ = http.Post(p.url, "application/json", bytes.NewBuffer(b))
	// if err != nil {
	// 	return model.Session{}, err
	// }
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return model.Session{}, err
	// }
	// r := model.Session{}
	// if err := json.Unmarshal(body, &r); err != nil {
	// 	return model.Session{}, err
	// }
	// return r, nil
}
func (p *paymentGateway) CancelCheckoutSession(ctx context.Context, sessionId string) (bool, error) {
	return true, nil
}
