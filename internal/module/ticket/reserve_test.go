package ticket

import (
	"context"
	"database/sql"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	paymentintegration "event_ticket/internal/platform/payment_integration"
	sticket "event_ticket/internal/storage/ticket"
	"net/http"
	"net/http/httptest"

	readtable "event_ticket/readTable"
	"fmt"
	"os"

	"testing"

	"github.com/cucumber/godog"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slog"
)

type contextKey string

type MockQueries struct {
	db.Querier
	Tkt db.Ticket
}

func (m *MockQueries) UpdateTicketStatus(ctx context.Context, arg db.UpdateTicketStatusParams) (db.Ticket, error) {
	m.Tkt = db.Ticket{
		TicketNo: arg.TicketNo,
		BusNo:    arg.BusNo,
		TripID:   arg.TripID,
		Status: sql.NullString{
			String: "Onhold",
			Valid:  true,
		}}

	return m.Tkt, nil
}

var mockqueries *MockQueries
var CallCount = 0

func TestReserveTicket(t *testing.T) {
	testCases := []struct {
		Name                string
		ScenarioInitializer func(sc *godog.ScenarioContext)
		FeatureFilepath     string
	}{
		{
			Name:                "user requests to reserve free ticket",
			ScenarioInitializer: reserveFreeticket,
			FeatureFilepath:     "reserve.feature",
		},
	}
	for _, tc := range testCases {

		t.Run(tc.Name, func(t *testing.T) {
			result := godog.TestSuite{
				Name:                 tc.Name,
				TestSuiteInitializer: nil,
				ScenarioInitializer:  tc.ScenarioInitializer,
				Options: &godog.Options{
					Paths:    []string{tc.FeatureFilepath},
					Format:   "pretty",
					TestingT: t,
				},
			}.Run()
			if result != 0 {
				t.Errorf("test failed")
			}
		})
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
	var dbQueriesKey contextKey = "key"
	ctx = context.WithValue(ctx, dbQueriesKey, mockqueries)

	return ctx, nil
}

func checkoutSessionRequestShouldBeSent(ctx context.Context) error {

	count := ctx.Value(contextKey("countKey")).(int)
	if count != 1 {
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
	mqueries, ok := ctx.Value(contextKey("key")).(*MockQueries)
	if !ok {
		return ctx, fmt.Errorf("no value found in context")
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := sticket.Init(logger, mqueries)

	url := ctx.Value(contextKey("serverURLKey")).(string)
	mpg := paymentintegration.Init(logger, url)

	moduleTicket := Init(slog.New(slog.NewJSONHandler(os.Stdout, nil)), store, mpg)
	_, err = moduleTicket.ReserveTicket(ctx, tickets[0].TicketNo, tickets[0].TripID, tickets[0].BusNo)
	if err != nil {
		return ctx, err
	}
	var countKey contextKey = "countKey"
	ctx = context.WithValue(ctx, countKey, CallCount)
	return ctx, nil
}

func reserveFreeticket(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CallCount++
		}))
		var serverKey contextKey = "serverURLKey"
		ctx = context.WithValue(ctx, serverKey, server.URL)
		return ctx, nil
	})
	sc.Step(`^a free ticket$`, aFreeTicket)
	sc.Step(`^user requests to reserve ticket$`, userRequestsToReserveTicket)
	sc.Step(`^the ticket status should be "([^"]*)"$`, theTicketStatusShouldBe)
	sc.Step(`^checkout session request should be sent$`, checkoutSessionRequestShouldBeSent)
}
