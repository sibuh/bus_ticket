package ticket

import (
	"context"
	"event_ticket/internal/model"
	readtable "event_ticket/readTable"
	"fmt"

	"testing"

	"github.com/cucumber/godog"
	"github.com/mitchellh/mapstructure"
)

type contextKey string

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
			ColumnType: readtable.String,
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
	var key contextKey
	ctx = context.WithValue(ctx, key, tickets[0])

	return ctx, nil
}

func checkoutSessionRequestShouldBeSent() error {

	return nil
}

func theTicketStatusShouldBe(arg1 string) error {
	return godog.ErrPending
}

func userRequestsToReserveTicket(arg1 *godog.Table, ctx context.Context) (context.Context, error) {
	_, ok := ctx.Value(contextKey("key")).(model.Ticket)
	if !ok {
		return ctx, fmt.Errorf("no value found in context")
	}

	// moduleTicket := Init(slog.New(slog.NewJSONHandler(os.Stdout, nil)), InitMock(tkt))

	return ctx, nil
}

func initScenario(sc *godog.ScenarioContext) {
	sc.Step(`^a free ticket$`, aFreeTicket)
	sc.Step(`^user requests to reserve ticket$`, userRequestsToReserveTicket)
	sc.Step(`^the ticket status should be "([^"]*)"$`, theTicketStatusShouldBe)
	sc.Step(`^checkout session request should be sent$`, checkoutSessionRequestShouldBeSent)
}
