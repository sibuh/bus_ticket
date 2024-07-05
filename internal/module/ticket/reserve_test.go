package ticket

import (
	"context"
	"event_ticket/internal/constant"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"fmt"
	"os"
	"time"

	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type reserveTicketTest struct {
	tkt         module.Ticket
	session     model.Session
	mockstorage *MockStorageTicket
	err         error
}

func TestReserveTicket(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := InitMock(model.Ticket{})
	platform := InitMockGateway(logger)
	reserveTkt := reserveTicketTest{
		tkt:         Init(logger, store, platform),
		mockstorage: store,
	}
	result := godog.TestSuite{
		Name:                 "ticket reservation test",
		TestSuiteInitializer: nil,
		ScenarioInitializer:  reserveTkt.InitializeScenario,
		Options: &godog.Options{
			Paths:  []string{"reserve.feature"},
			Format: "pretty",
		},
	}.Run()
	if result != 0 {
		t.Errorf("test failed")
	}

}
func (r *reserveTicketTest) ticketMustBeSetForSale(status string) error {
	if r.mockstorage.Tkt.Status != status {
		return fmt.Errorf("ticket status is not set free")
	}
	return nil
}
func (r *reserveTicketTest) ticketReservstionDoNotSucceedWithInSDuration(delay int) error {
	time.Sleep(10 * time.Second)
	return nil
}
func (r *reserveTicketTest) userShouldGetErrorMessage(errMsg string) error {
	if r.err.Error() != errMsg {
		return fmt.Errorf("want:%s got: %s", errMsg, r.err.Error())
	}
	return nil
}

func (r *reserveTicketTest) theTicketStatusShouldBe(status string) error {
	if r.mockstorage.Tkt.Status != status {
		return fmt.Errorf("want: %s got: %s", status, r.mockstorage.Tkt.Status)
	}
	return nil
}

func (r *reserveTicketTest) theUserShouldGetCheckoutUrl() error {
	if r.session.PaymentUrl == "" {
		return fmt.Errorf("checkout url not returned")
	}
	return nil
}

func (r *reserveTicketTest) ticketNumberOfBusNumberForTripOfIdIs(tktNo, busNo, tripId int, status string) error {
	r.mockstorage.Tkt = model.Ticket{
		TripId:   int32(tripId),
		TicketNo: int32(tktNo),
		BusNo:    int32(busNo),
		Status:   status,
	}

	return nil
}

func (r *reserveTicketTest) userRequestsToReserveTicketNumberOfTrip(tktNo, tripId int) error {

	r.session, r.err = r.tkt.ReserveTicket(context.Background(), int32(tktNo), int32(tripId))

	return nil
}
func (r *reserveTicketTest) cancelCheckoutSessionIsSentToPaymentGateway() error {
	return nil
}

func (r *reserveTicketTest) checkPaymentStatusOnPaymentGateway() error {
	return nil
}

func (r *reserveTicketTest) checkoutSessionIsCreated(arg1 *godog.Table) error {
	r.session = model.Session{
		ID: uuid.NewString(),
		Tkt: model.Ticket{
			TripId:   int32(778),
			TicketNo: int32(12),
			BusNo:    int32(10),
			Status:   string(constant.Onhold),
		},
		CreatedAt:     time.Now(),
		PaymentStatus: string(constant.Pending),
	}
	return nil
}

func (r *reserveTicketTest) paymentCancelationResponseIsSuccessful() error {
	return nil
}

func (r *reserveTicketTest) paymentStatusCheckoutSessionReturnsForCheckoutSession(arg1 string) error {
	return nil
}

func (r *reserveTicketTest) paymentStatusForCheckoutSessionReturns(arg1 string) error {
	return nil
}

func (r *reserveTicketTest) paymentStatusIsRequestedForCheckoutSession() error {
	return nil
}

func (r *reserveTicketTest) ticketMustBeSetToStatus(arg1 string) error {
	return nil
}

func (r *reserveTicketTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^user should get error message "([^"]*)"$`, r.userShouldGetErrorMessage)
	ctx.Step(`^the ticket status should be "([^"]*)"$`, r.theTicketStatusShouldBe)
	ctx.Step(`^the user should get checkout url$`, r.theUserShouldGetCheckoutUrl)
	ctx.Step(`^ticket number (\d+) of bus number (\d+) for trip of id (\d+) is "([^"]*)"$`, r.ticketNumberOfBusNumberForTripOfIdIs)
	ctx.Step(`^user requests to reserve ticket number (\d+) of trip (\d+)$`, r.userRequestsToReserveTicketNumberOfTrip)
	ctx.Step(`^ticket must be set "([^"]*)" for sale$`, r.ticketMustBeSetForSale)
	ctx.Step(`^ticket reservstion do not succeed with in (\d+)s duration$`, r.ticketReservstionDoNotSucceedWithInSDuration)
	ctx.Step(`^cancel checkout session is sent to payment gateway$`, r.cancelCheckoutSessionIsSentToPaymentGateway)
	ctx.Step(`^check payment status on payment gateway$`, r.checkPaymentStatusOnPaymentGateway)
	ctx.Step(`^checkout session is created$`, r.checkoutSessionIsCreated)
	ctx.Step(`^payment cancelation response is successful$`, r.paymentCancelationResponseIsSuccessful)
	ctx.Step(`^payment status checkout session returns "([^"]*)" for checkout session$`, r.paymentStatusCheckoutSessionReturnsForCheckoutSession)
	ctx.Step(`^payment status for checkout session returns "([^"]*)"$`, r.paymentStatusForCheckoutSessionReturns)
	ctx.Step(`^payment status is requested for checkout session$`, r.paymentStatusIsRequestedForCheckoutSession)
	ctx.Step(`^ticket must be set to "([^"]*)" status$`, r.ticketMustBeSetToStatus)
}
