package ticket

import (
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	storageTkt "event_ticket/internal/storage/ticket"
	"fmt"

	"testing"

	"github.com/cucumber/godog"
	"golang.org/x/exp/slog"
)

type holdTicketTest struct {
	tkt         module.Ticket
	mockstorage *storageTkt.MockStorageTicket
}

func TestHoldTicket(t *testing.T) {
	logger := slog.Logger{}
	store := storageTkt.InitMock(model.Ticket{})
	holdTkt := holdTicketTest{
		tkt:         Init(logger, store),
		mockstorage: store,
	}
	result := godog.TestSuite{
		Name:                 "test hold ticket",
		TestSuiteInitializer: nil,
		ScenarioInitializer:  holdTkt.InitializeScenario,
		Options: &godog.Options{
			Paths:  []string{"reserve.feature"},
			Format: "pretty",
		},
	}.Run()
	if result != 0 {
		t.Errorf("test failed")
	}

}
func (h *holdTicketTest) theTicketStatusShouldBe(status string) error {
	if h.mockstorage.Tkt.Status != status {
		return fmt.Errorf("want: %s got: %s", status, h.mockstorage.Tkt.Status)
	}

	return nil
}

func (h *holdTicketTest) theUserShouldGetCkeckoutUrl() error {
	return nil
}

func (h *holdTicketTest) ticketNumberOfBusNumberForTripOfIdIs(tktNo, busNo, tripId int, status string) error {
	_, err := h.mockstorage.AddTicket(int32(tktNo), int32(busNo), int32(tripId), status)
	fmt.Println("mock storage:", h.mockstorage)
	if err != nil {
		return err
	}

	return nil
}

func (h *holdTicketTest) userRequestsToReserveTicketNumberOfTrip(tktNo, tripId int) error {

	tkt, err := h.tkt.HoldTicket(int32(tktNo), int32(tripId))
	if err != nil {
		return err
	}
	fmt.Println("tkt:--->", tkt)

	return nil
}

func (h *holdTicketTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the ticket status should be "([^"]*)"$`, h.theTicketStatusShouldBe)
	ctx.Step(`^the user should get ckeckout url$`, h.theUserShouldGetCkeckoutUrl)
	ctx.Step(`^ticket number (\d+) of bus number (\d+) for trip of id (\d+) is "([^"]*)"$`, h.ticketNumberOfBusNumberForTripOfIdIs)
	ctx.Step(`^user requests to reserve ticket number (\d+) of trip (\d+)$`, h.userRequestsToReserveTicketNumberOfTrip)
}
