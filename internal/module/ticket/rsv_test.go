package ticket

import (
	"context"
	"database/sql"
	"encoding/json"
	"event_ticket/internal/constant"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	sticket "event_ticket/internal/storage/ticket"
	"net/http"
	"net/http/httptest"
	"time"

	readtable "event_ticket/readTable"
	"fmt"
	"os"

	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slog"
)

type contextKey string

var mockqueries *MockQueries

func TestReserveTicket(t *testing.T) {
	result := godog.TestSuite{
		Name:                 "ticket reservation test",
		TestSuiteInitializer: nil,
		ScenarioInitializer:  initScenario,
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

func aFreeTicket(ctx context.Context, t *godog.Table) (context.Context, error) {
	result, err := readtable.ReadTableData(t, []readtable.Column{
		{
			ColimnName: "ticket_no",
			ColumnType: readtable.Int,
		},
		{
			ColimnName: "bus_no",
			ColumnType: readtable.Int,
		},
		{
			ColimnName: "trip_id",
			ColumnType: readtable.Int,
		},
		{
			ColimnName: "status",
			ColumnType: readtable.String,
		},
	})
	if err != nil {
		return ctx, err
	}

	var tickets []model.Ticket

	err = mapstructure.Decode(result, &tickets)
	if err != nil {
		return ctx, err
	}
	mockqueries = &MockQueries{
		Tkt: db.Ticket{
			TripID: tickets[0].TripID,
			BusNo:  tickets[0].TicketNo,
			Status: sql.NullString{
				String: tickets[0].Status,
				Valid:  true,
			},
			TicketNo: tickets[0].TicketNo,
		},
	}
	var dbQueriesKey contextKey
	ctx = context.WithValue(ctx, dbQueriesKey, mockqueries)
	return ctx, nil
}

func checkoutSessionRequestShouldBeSent(ctx context.Context) error {
	session := ctx.Value(contextKey("sessionKey")).(model.Session)
	if session.PaymentUrl == "" {
		return fmt.Errorf("checkout session not created")
	}
	return nil
}

func theTicketStatusShouldBe(arg1 string) error {
	if mockqueries.Tkt.Status.String != arg1 {
		return fmt.Errorf("tickets status not updated want %s: got: %s", arg1, mockqueries.Tkt.Status.String)
	}
	return nil
}

func userRequestsToReserveTicket(ctx context.Context, arg1 *godog.Table) (context.Context, error) {
	result, err := readtable.ReadTableData(arg1, []readtable.Column{
		{
			ColimnName: "ticket_no",
			ColumnType: readtable.Int,
		},
		{
			ColimnName: "bus_no",
			ColumnType: readtable.Int,
		},
		{
			ColimnName: "trip_id",
			ColumnType: readtable.Int,
		},
	})
	if err != nil {
		return ctx, err
	}

	var tickets []model.Ticket

	err = mapstructure.Decode(result, &tickets)
	if err != nil {
		return ctx, err
	}
	mqueries, ok := ctx.Value(contextKey("dbQueriesKey")).(*MockQueries)
	if !ok {
		return ctx, fmt.Errorf("no value found in context")
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := sticket.Init(logger, mqueries)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := model.Session{
			ID: uuid.NewString(),
			Tkt: model.Ticket{
				TripID:   778,
				BusNo:    10,
				TicketNo: 12,
				Status:   string(constant.Onhold),
			},
			PaymentStatus: string(constant.Pending),
			PaymentUrl:    "http://payment.com/sessionId",
			TotalAmount:   200,
			CreatedAt:     time.Now(),
		}
		b, _ := json.Marshal(sess)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}))
	mpg := InitMockGateway(logger, server.URL)

	moduleTicket := Init(slog.New(slog.NewJSONHandler(os.Stdout, nil)), store, mpg)
	session, err := moduleTicket.ReserveTicket(ctx, tickets[0].TicketNo, tickets[0].TripID, tickets[0].BusNo)
	if err != nil {
		return ctx, err
	}
	var sessionKey contextKey
	ctx = context.WithValue(ctx, sessionKey, session)
	return ctx, nil
}

func initScenario(sc *godog.ScenarioContext) {
	sc.Step(`^a free ticket$`, aFreeTicket)
	sc.Step(`^user requests to reserve ticket$`, userRequestsToReserveTicket)
	sc.Step(`^the ticket status should be "([^"]*)"$`, theTicketStatusShouldBe)
	sc.Step(`^checkout session request should be sent$`, checkoutSessionRequestShouldBeSent)
}
