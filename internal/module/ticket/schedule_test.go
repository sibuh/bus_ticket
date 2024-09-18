package ticket

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
)

// type Mqueries struct {
// 	Ssn db.Session
// 	Tkt db.Ticket
// 	db.Querier
// }

// func (mq *Mqueries) GetTicketStatus(ctx context.Context, sid string) (string, error) {
// 	return mq.Tkt.Status, nil
// }

// var callCount int
// var ch = make(chan string)
// var id string

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

// func noPaymentStatusCheckRequestShouldBeSentWithinS(ctx context.Context, arg1 int) error {
// 	t := time.NewTimer(time.Duration(arg1 - 1))
// 	select {
// 	case <-t.C:
// 		return nil
// 	case <-channel:
// 		return fmt.Errorf("payment status check request should not be sent to gateway before %d secs", arg1)
// 	}

// }

func paymentStatusCheckRequestIsScheduledForCheckoutSession(ctx context.Context) (context.Context, error) {

	id := uuid.NewString()
	var channel = make(chan string, 1)

	callCount := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount = callCount + "1"
		fmt.Println("channel:", callCount)
	}))

	go Scheduler(id, channel, 2, func() error {
		_, err := http.Get(server.URL)
		if err != nil {
			return err
		}
		return nil
	})
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
	// select {
	// case <-t.C:
	// 	return fmt.Errorf("payment status check request not sent to gateway")
	// case <-channel:
	// 	return nil
	// }

}

// func scheduledProcessShouldBeTerminated() error {
// 	t := time.NewTimer(2 * time.Second)
// 	select {
// 	case <-channel:
// 		return fmt.Errorf("check payment status request must be cancelled")
// 	case <-t.C:
// 		return nil
// 	}

// }

// func successOrFailureCallbackArrivesForCheckoutSession(ctx context.Context) error {
// 	id, ok := ctx.Value(contextKey("id-key")).(string)
// 	if !ok {
// 		return fmt.Errorf("failed to get id from context")
// 	}
// 	ch <- id
// 	return nil
// }

func InitializeScenario(ctx *godog.ScenarioContext) {
	// ctx.Step(`^no payment status check request should be sent within (\d+)s$`,
	// 	noPaymentStatusCheckRequestShouldBeSentWithinS)
	ctx.Step(`^payment status check request is scheduled for checkout session$`,
		paymentStatusCheckRequestIsScheduledForCheckoutSession)
	ctx.Step(`^payment status check request should be sent to payment gateway after (\d+)s$`,
		paymentStatusCheckRequestShouldBeSentToPaymentGatewayAfterS)
	// ctx.Step(`^scheduled process should be terminated$`, scheduledProcessShouldBeTerminated)
	// ctx.Step(`^success or failure callback arrives for checkout session$`,
	// 	successOrFailureCallbackArrivesForCheckoutSession)

}
