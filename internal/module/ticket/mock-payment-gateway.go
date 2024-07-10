package ticket

import (
	"context"
	"event_ticket/internal/constant"
	"event_ticket/internal/model"
	"event_ticket/internal/platform"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type MockPaymentGateWay struct {
	logger *slog.Logger
	url    string
	platform.PaymentGatewayIntegrator
}

func InitMockGateway(logger *slog.Logger, url string) *MockPaymentGateWay {
	return &MockPaymentGateWay{logger: logger, url: url}
}

func (m *MockPaymentGateWay) CreateCheckoutSession(ctx context.Context, ticketInfo model.Ticket) (model.Session, error) {
	if ticketInfo.TripId == int32(779) {
		return model.Session{}, fmt.Errorf("failed to create checkout session")
	}
	return model.Session{
		ID:            uuid.NewString(),
		Tkt:           model.Ticket{},
		PaymentStatus: string(constant.Pending),
		PaymentUrl:    "https://chapa.com/checkout/session-id",
		CreatedAt:     time.Now(),
	}, nil
}
