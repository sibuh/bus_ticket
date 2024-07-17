package ticket

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
)

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

func checkoutSessionIsSuccessfullyCreated(ctx context.Context) error {
	// sticket.Init(logger, mockqueries)
	return nil
}

func paymentDoNotCompleteWithinSeconds(arg1 int) error {
	return nil
}

func paymentStatusCheckRequestShouldBeSentToGateway() error {
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^checkout session is successfully created$`, checkoutSessionIsSuccessfullyCreated)
	ctx.Step(`^payment do not complete within (\d+) seconds$`, paymentDoNotCompleteWithinSeconds)
	ctx.Step(`^payment status check request should be sent to gateway$`, paymentStatusCheckRequestShouldBeSentToGateway)
}
