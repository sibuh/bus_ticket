package ticket

import (
	"context"
	"encoding/json"
	"event_ticket/internal/constant"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	paymentintegration "event_ticket/internal/platform/payment_integration"
	"event_ticket/internal/storage/session"
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
	Ssn db.Session
}

func (m *MockQueries) UpdateTicketStatus(ctx context.Context, arg db.UpdateTicketStatusParams) (db.Ticket, error) {
	m.Tkt.Status = arg.Status

	return m.Tkt, nil
}
func (m *MockQueries) GetTicket(ctx context.Context, id string) (db.Ticket, error) {
	return m.Tkt, nil
}
func (m *MockQueries) StoreCheckoutSession(ctx context.Context, arg db.StoreCheckoutSessionParams) (db.Session, error) {
	m.Ssn = db.Session{
		ID:            arg.ID,
		TicketID:      arg.TicketID,
		PaymentStatus: arg.PaymentStatus,
		PaymentURL:    arg.PaymentURL,
		CancelURl:     arg.CancelURL,
		Amount:        arg.Amount,
		CreatedAt:     arg.CreatedAt,
	}
	return m.Ssn, nil
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func TestReserveTicket(t *testing.T) {
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
}

func aTicket(ctx context.Context, arg1 string) (context.Context, error) {
	mockqueries := &MockQueries{
		Tkt: db.Ticket{
			ID:       uuid.NewString(),
			TripID:   21,
			BusNo:    2321,
			Status:   arg1,
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

func theTicketStatusShouldBe(ctx context.Context, arg1 string) error {
	queries, ok := ctx.Value(contextKey("ticket-data")).(*MockQueries)
	if !ok {
		return fmt.Errorf("failed to get ticket from context")
	}
	if queries.Tkt.Status != arg1 {
		return fmt.Errorf("ticket status not updated want %s: got: %s", arg1, queries.Tkt.Status)
	}
	return nil
}

func userRequestsToReserveTicket(ctx context.Context) (context.Context, error) {
	var CallCount = 0

	mqueries, ok := ctx.Value(contextKey("ticket-data")).(*MockQueries)
	if !ok {
		return ctx, fmt.Errorf("no value found in context")
	}

	store := sticket.Init(logger, mqueries)

	url := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CallCount++
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(model.Session{})

		w.Write(b)
	})).URL
	mpg := paymentintegration.Init(logger, url)
	ssn := session.Init(logger, mqueries)
	moduleTicket := Init(logger, store, mpg, ssn)
	_, err := moduleTicket.ReserveTicket(ctx, model.ReserveTicketRequest{ID: mqueries.Tkt.ID}, func() {})
	if err != nil {
		var errorKey contextKey = "error-key"
		ctx = context.WithValue(ctx, errorKey, err)
		return ctx, nil
	}
	var countKey contextKey = "count-key"
	ctx = context.WithValue(ctx, countKey, CallCount)
	return ctx, nil
}

func checkoutSessionShouldBeStored(ctx context.Context) error {
	session, ok := ctx.Value(contextKey("session-key")).(model.Session)
	if !ok {
		return fmt.Errorf("could not find session data from context")
	}
	if session.PaymentURL == "" {
		return fmt.Errorf("payment url is empty")
	}
	queries, ok := ctx.Value(contextKey("ticket-data")).(*MockQueries)
	if !ok {
		return fmt.Errorf("failed to ticket data from context")
	}
	if queries.Ssn.PaymentURL != session.PaymentURL {
		return fmt.Errorf("paymentURL not updated want:%s got:%s", session.PaymentURL, queries.Ssn.PaymentURL)
	}
	return nil
}

func createCheckoutSessionSucceedsForReservingTicketRequest(ctx context.Context) (context.Context, error) {
	queries := ctx.Value(contextKey("ticket-data")).(*MockQueries)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	store := sticket.Init(logger, queries)
	ssn := session.Init(logger, queries)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(model.Session{
			ID:            uuid.NewString(),
			TicketID:      uuid.NewString(),
			PaymentStatus: string(constant.Pending),
			PaymentURL:    "http://payment/session_id",
			CancelURL:     "http://cancel_url",
			Amount:        400,
			CreatedAt:     time.Now(),
		})

		w.Write(b)
	}))
	url := server.URL
	platform := paymentintegration.Init(logger, url)
	mod := Init(logger, store, platform, ssn)
	var SchedulerCount int = 0

	session, err := mod.ReserveTicket(ctx, model.ReserveTicketRequest{ID: queries.Tkt.ID}, func() { SchedulerCount++ })
	if err != nil {
		return nil, err
	}
	var sessionKey contextKey = "session-key"
	ctx = context.WithValue(ctx, sessionKey, session)
	var scheduleKey contextKey = "schedule-key"
	ctx = context.WithValue(ctx, scheduleKey, SchedulerCount)
	return ctx, nil
}

func theUserShouldGetCheckoutUrl(ctx context.Context) error {
	session := ctx.Value(contextKey("session-key")).(model.Session)
	if session.PaymentURL == "" {
		return fmt.Errorf("no payment url is returned ")
	}

	return nil
}
func checkoutSessionCreationFailsDuringReserveTicketRequest(ctx context.Context) (context.Context, error) {
	queries := ctx.Value(contextKey("ticket-data")).(*MockQueries)
	store := sticket.Init(logger, queries)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	platform := paymentintegration.Init(logger, server.URL)
	session := session.Init(logger, queries)
	mod := Init(logger, store, platform, session)

	_, err := mod.ReserveTicket(ctx, model.ReserveTicketRequest{ID: queries.Tkt.ID, Status: string(constant.Onhold)}, func() {
	})
	if err == nil {
		return ctx, fmt.Errorf("expected non-nil error but did not get any error")
	}
	var errorKey contextKey = "error-key"
	ctx = context.WithValue(ctx, errorKey, err)
	return ctx, nil
}

func userShouldGetErrorMessage(ctx context.Context, arg1 string) error {
	err := ctx.Value(contextKey("error-key")).(error)
	if arg1 != err.Error() {
		return fmt.Errorf("expected: %s got: %s", arg1, err.Error())
	}

	return nil
}
func onholdtimeoutProcessShouldBeScheduled(ctx context.Context) error {
	count, ok := ctx.Value(contextKey("schedule-key")).(int)
	if !ok {
		return fmt.Errorf("failed to get value from context")
	}
	if count != 1 {
		return fmt.Errorf("scheduler not called")
	}

	return nil
}

func ReserveFreeticketScenario(sc *godog.ScenarioContext) {
	sc.Step(`^a "([^"]*)" ticket$`, aTicket)
	sc.Step(`^user requests to reserve ticket$`, userRequestsToReserveTicket)
	sc.Step(`^the ticket status should be "([^"]*)"$`, theTicketStatusShouldBe)
	sc.Step(`^checkout session request should be sent$`, checkoutSessionRequestShouldBeSent)
	sc.Step(`^checkout session should be stored$`, checkoutSessionShouldBeStored)
	sc.Step(`^create checkout session succeeds for reserving ticket request$`,
		createCheckoutSessionSucceedsForReservingTicketRequest)
	sc.Step(`^the user should get checkout url$`, theUserShouldGetCheckoutUrl)
	sc.Step(`^checkout session creation fails during reserve ticket request$`,
		checkoutSessionCreationFailsDuringReserveTicketRequest)
	sc.Step(`^user should get error message "([^"]*)"$`, userShouldGetErrorMessage)
	sc.Step(`^onhold-timeout process should be scheduled$`, onholdtimeoutProcessShouldBeScheduled)
}
