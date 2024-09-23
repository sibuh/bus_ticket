package ticket

import (
	"context"
	"event_ticket/internal/model"
	"event_ticket/internal/module/callback"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
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

func noPaymentStatusCheckRequestShouldBeSentWithinS(ctx context.Context, arg1 int) error {
	time.Sleep(time.Duration(arg1-1) * time.Second)
	c := ctx.Value(contextKey("count")).(*string)
	if *c != "" {
		fmt.Printf("count value, %v", c)
		return fmt.Errorf("payment status check request should not be sent to gateway")
	}
	return nil
}

func paymentStatusCheckRequestIsScheduledForCheckoutSession(ctx context.Context) (context.Context, error) {

	id := uuid.NewString()
	var channel = make(chan string, 1)

	callCount := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount = callCount + "1"
		fmt.Println("channel:", callCount)
	}))

	// scheduler registers callback
	cback := callback.Init()

	cback.Scheduler.Append(id, channel)
	go cback.Scheduler.Schedule(id, channel, 2, func() error {
		_, err := http.Get(server.URL)
		if err != nil {
			return err
		}
		return nil
	})
	ctx = context.WithValue(ctx, contextKey("callback"), cback)
	ctx = context.WithValue(ctx, contextKey("sessionId"), id)
	ctx = context.WithValue(ctx, contextKey("count"), &callCount)

	return ctx, nil
}

func paymentStatusCheckRequestShouldBeSentToPaymentGatewayAfterS(ctx context.Context, arg1 int) error {
	time.Sleep(time.Duration(arg1+1) * time.Second)
	c := ctx.Value(contextKey("count")).(*string)
	if *c != "1" {
		fmt.Printf("count value, %v", c)
		return fmt.Errorf("payment status check request not sent to gateway")
	}
	return nil
}

func scheduledProcessShouldBeTerminated(ctx context.Context) error {
	cback := ctx.Value(contextKey("callback")).(*callback.Callback)
	sessionId := ctx.Value(contextKey("sessionId")).(string)

	ch := cback.Scheduler.Get(sessionId)

	if ch != nil {
		return fmt.Errorf("Scheduled process should have been removed")
	}
	return nil
}

func successOrFailureCallbackArrivesForCheckoutSession(ctx context.Context) context.Context {
	// callback module initiate
	cback := ctx.Value(contextKey("callback")).(*callback.Callback)
	sessionId := ctx.Value(contextKey("sessionId")).(string)

	samplePayload := model.Payment{
		IntentID: sessionId,
	}
	cback.HandlePaymentStatusUpdate(samplePayload)
	return ctx
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^no payment status check request should be sent within (\d+)s$`,
		noPaymentStatusCheckRequestShouldBeSentWithinS)
	ctx.Step(`^payment status check request is scheduled for checkout session$`,
		paymentStatusCheckRequestIsScheduledForCheckoutSession)
	ctx.Step(`^payment status check request should be sent to payment gateway after (\d+)s$`,
		paymentStatusCheckRequestShouldBeSentToPaymentGatewayAfterS)
	ctx.Step(`^scheduled process should be terminated$`, scheduledProcessShouldBeTerminated)
	ctx.Step(`^success or failure callback arrives for checkout session$`,
		successOrFailureCallbackArrivesForCheckoutSession)

}
