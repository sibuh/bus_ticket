package ticket

import (
	"context"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	paymentintegration "event_ticket/internal/platform/payment_integration"
	storageTkt "event_ticket/internal/storage/ticket"
	"fmt"
	"time"

	"testing"

	"github.com/cucumber/godog"
	"golang.org/x/exp/slog"
)

type reserveTicketTest struct {
	tkt         module.Ticket
	mockstorage *storageTkt.MockStorageTicket
	url         string
	err         error
}

type MockStorageTicket struct {
	Tkt model.Ticket
}

func InitMock(tkt model.Ticket) *MockStorageTicket {
	return &MockStorageTicket{Tkt: tkt}
}
func (m *MockStorageTicket) ReserveTicket(ticketNo, tripId int32) (model.Ticket, error) {

	m.Tkt.Status = "Onhold"
	return m.Tkt, nil

}
func (m *MockStorageTicket) GetTicket(tktNo, tripId int32) (model.Ticket, error) {
	return m.Tkt, nil
}
func (m *MockStorageTicket) UnholdTicket(tktNo, tripId int32) (model.Ticket, error) {
	if m.Tkt.Status == "Onhold" {
		m.Tkt.Status = "Free"
		return m.Tkt, nil
	}
	return model.Ticket{}, fmt.Errorf("failed to unhold ticket")
}

func TestReserveTicket(t *testing.T) {
	logger := slog.Logger{}
	store := storageTkt.InitMock(model.Ticket{})
	platform := paymentintegration.InitMock(logger)
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
	if r.url == "" {
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

	r.url, r.err = r.tkt.ReserveTicket(context.Background(), int32(tktNo), int32(tripId))

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
}
