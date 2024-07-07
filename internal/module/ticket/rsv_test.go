package ticket

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
)

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
	fmt.Printf("godog table, %v", t)
	return ctx, nil
}

func checkoutSessionRequestShouldBeSent() error {
	return godog.ErrPending
}

func theTicketStatusShouldBe(arg1 string) error {
	return godog.ErrPending
}

func userRequestsToReserveTicket() error {
	return godog.ErrPending
}

func initScenario(sc *godog.ScenarioContext) {
	sc.Step(`^a free ticket$`, aFreeTicket)
	sc.Step(`^user requests to reserve ticket$`, userRequestsToReserveTicket)
	sc.Step(`^the ticket status should be "([^"]*)"$`, theTicketStatusShouldBe)
	sc.Step(`^checkout session request should be sent$`, checkoutSessionRequestShouldBeSent)
}
