package ticket

import (
	"context"
	"event_ticket/internal/constant"
	"event_ticket/internal/data/db"
	"event_ticket/internal/storage/session"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type Mqueries struct {
	Ssn db.Session
	Tkt db.Ticket
	db.Querier
}

func (mq *Mqueries) GetTicketStatus(ctx context.Context, sid string) (string, error) {
	return mq.Tkt.Status, nil
}

func TestScheduleOntimeoutProcess(t *testing.T) {
	result := godog.TestSuite{
		Name:                 "schedule ontimeout process test",
		TestSuiteInitializer: nil,
		ScenarioInitializer:  InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"schedule.feature"},
			TestingT: t,
		},
	}.Run()
	if result != 0 {
		t.Errorf("schedule ontimeout process failed")
	}
}

func checkoutSessionIsSuccessfullyCreated(ctx context.Context) (context.Context, error) {
	queries := &Mqueries{
		Ssn: db.Session{
			ID:            uuid.NewString(),
			TicketID:      uuid.NewString(),
			PaymentStatus: string(constant.Onhold),
			PaymentURL:    "http://paymet.com/session_id",
			CancelURl:     "http://payment.com/cancel",
			Amount:        455,
			CreatedAt:     time.Now(),
		},
		Tkt: db.Ticket{
			ID:       uuid.NewString(),
			TripID:   45,
			BusNo:    45,
			TicketNo: 45,
			Status:   string(constant.Onhold),
		},
	}
	var queriesKey contextKey = "data-key"
	ctx = context.WithValue(ctx, queriesKey, queries)
	return ctx, nil
}

func paymentDoNotCompleteWithinSeconds(arg1 int) error {
	time.Sleep(time.Second)
	return nil
}

func paymentStatusCheckRequestShouldBeSentToGateway(ctx context.Context) error {
	queries, ok := ctx.Value(contextKey("data-key")).(*Mqueries)
	if !ok {
		return fmt.Errorf("failed to get ticket data from context")
	}
	callCount := 0
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
	}))
	session := session.Init(logger, queries)

	Scheduler(session, server.URL, queries.Ssn.ID, logger)

	if callCount != 1 {
		return fmt.Errorf("payment status check request not sent")
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^checkout session is successfully created$`, checkoutSessionIsSuccessfullyCreated)
	ctx.Step(`^payment do not complete within (\d+) seconds$`, paymentDoNotCompleteWithinSeconds)
	ctx.Step(`^payment status check request should be sent to gateway$`, paymentStatusCheckRequestShouldBeSentToGateway)
}
