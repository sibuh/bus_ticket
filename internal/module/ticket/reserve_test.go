package ticket

import (
	"context"
	"encoding/json"
	"event_ticket/internal/constant"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	paymentintegration "event_ticket/internal/platform/payment_integration"
	sticket "event_ticket/internal/storage/ticket"
	"net/http"
	"net/http/httptest"
	"time"

	"fmt"
	"os"

	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type contextKey string

type MockQueries struct {
	db.Querier
	Tkt db.Ticket
}

func (m *MockQueries) UpdateTicketStatus(ctx context.Context, arg db.UpdateTicketStatusParams) (db.Ticket, error) {
	m.Tkt.Status = arg.Status

	return m.Tkt, nil
}
func (m *MockQueries) GetTicket(ctx context.Context, id string) (db.Ticket, error) {
	return m.Tkt, nil
}

var mockqueries *MockQueries
var CallCount = 0

func TestReserveTicket(t *testing.T) {
	// testCases := []struct {
	// 	Name                string
	// 	ScenarioInitializer func(sc *godog.ScenarioContext)
	// 	FeatureFilepath     string
	// }{
	// 	{
	// 		Name:                "user requests to reserve free ticket",
	// 		ScenarioInitializer: ReserveFreeticketScenario,
	// 		FeatureFilepath:     "reserve.feature",
	// 	},
	// 	{
	// 		Name:                "checkout session is created successfully",
	// 		ScenarioInitializer: CheckoutSessionSuccessScenario,
	// 		FeatureFilepath:     "checkout_session_created.feature",
	// 	},
	// }
	// for _, tc := range testCases {
	// 	t.Run(tc.Name, func(t *testing.T) {
	result := godog.TestSuite{
		Name:                 "reserve ticket",
		TestSuiteInitializer: nil,
		ScenarioInitializer:  ReserveFreeticketScenario,
		Options: &godog.Options{
			Paths:    []string{"reserve.feature"},
			Format:   "pretty",
			TestingT: t,
		},
	}.Run()
	if result != 0 {
		t.Errorf("test failed")
	}
	// 	})
	// }

}

func aFreeTicket(ctx context.Context) (context.Context, error) {
	mockqueries = &MockQueries{
		Tkt: db.Ticket{
			ID:       uuid.NewString(),
			TripID:   21,
			BusNo:    2321,
			Status:   string(constant.Free),
			TicketNo: 23,
		},
	}

	var dbQueriesKey contextKey = "ticket-data"
	ctx = context.WithValue(ctx, dbQueriesKey, mockqueries)

	return ctx, nil
}

func checkoutSessionRequestShouldBeSent(ctx context.Context) error {

	count := ctx.Value(contextKey("count-key")).(int)
	if count != 1 {
		return fmt.Errorf("checkout session not created")
	}
	return nil
}

func theTicketStatusShouldBe(arg1 string) error {
	if mockqueries.Tkt.Status != arg1 {
		return fmt.Errorf("ticket status not updated want %s: got: %s", arg1, mockqueries.Tkt.Status)
	}
	return nil
}

func userRequestsToReserveTicket(ctx context.Context) (context.Context, error) {

	mqueries, ok := ctx.Value(contextKey("ticket-data")).(*MockQueries)
	if !ok {
		return ctx, fmt.Errorf("no value found in context")
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := sticket.Init(logger, mqueries)

	url := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CallCount++
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(model.Session{})

		w.Write(b)
	})).URL
	mpg := paymentintegration.Init(logger, url)

	moduleTicket := Init(logger, store, mpg)
	_, err := moduleTicket.ReserveTicket(ctx, model.ReserveTicketRequest{ID: mqueries.Tkt.ID})
	if err != nil {
		return ctx, err
	}
	var countKey contextKey = "count-key"
	ctx = context.WithValue(ctx, countKey, CallCount)
	return ctx, nil
}

func checkoutSessionShouldBeStored(ctx context.Context) error {
	fmt.Println("session not stored yet")
	return nil
}

func createCheckoutSessionSucceedsForReservingTicketRequest(ctx context.Context) (context.Context, error) {
	queries := ctx.Value(contextKey("ticket-data")).(*MockQueries)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	store := sticket.Init(logger, queries)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(model.Session{
			ID: uuid.NewString(),
			Tkt: model.Ticket{
				ID:       uuid.NewString(),
				BusNo:    123,
				TicketNo: 456,
				TripID:   324,
				Status:   string(constant.Onhold),
			},
			PaymentStatus: string(constant.Pending),
			PaymentUrl:    "http://payment/session_id",
			CancelUrl:     "http://cancel_url",
			TotalAmount:   400,
			CreatedAt:     time.Now(),
		})

		w.Write(b)
	}))
	url := server.URL
	platform := paymentintegration.Init(logger, url)
	mod := Init(logger, store, platform)
	session, err := mod.ReserveTicket(ctx, model.ReserveTicketRequest{ID: queries.Tkt.ID})
	if err != nil {
		return nil, err
	}
	var sessionKey contextKey = "session-key"
	ctx = context.WithValue(ctx, sessionKey, session)
	return ctx, nil
}

func theUserShouldGetCheckoutUrl(ctx context.Context) error {
	session := ctx.Value(contextKey("session-key")).(model.Session)
	if session.PaymentUrl == "" {
		return fmt.Errorf("no payment url is returned ")
	}

	return nil
}

// func CheckoutSessionSuccessScenario(ctx *godog.ScenarioContext) {
// 	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {

// 		var serverKey contextKey = "server-url-key"
// 		ctx = context.WithValue(ctx, serverKey, server.URL)
// 		return ctx, nil
// 	})

// }

func ReserveFreeticketScenario(sc *godog.ScenarioContext) {
	// sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {

	// 	var serverKey contextKey = "server-url-key"
	// 	ctx = context.WithValue(ctx, serverKey, server.URL)
	// 	return ctx, nil
	// })
	sc.Step(`^a free ticket$`, aFreeTicket)
	sc.Step(`^user requests to reserve ticket$`, userRequestsToReserveTicket)
	sc.Step(`^the ticket status should be "([^"]*)"$`, theTicketStatusShouldBe)
	sc.Step(`^checkout session request should be sent$`, checkoutSessionRequestShouldBeSent)
	sc.Step(`^a free ticket$`, aFreeTicket)
	sc.Step(`^checkout session should be stored$`, checkoutSessionShouldBeStored)
	sc.Step(`^create checkout session succeeds for reserving ticket request$`, createCheckoutSessionSucceedsForReservingTicketRequest)
	sc.Step(`^the user should get checkout url$`, theUserShouldGetCheckoutUrl)
}
