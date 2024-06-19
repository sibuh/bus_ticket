package ticket_test

import (
	"event_ticket/internal/module/ticket"
	"testing"

	"github.com/cucumber/godog"
)

func TestHoldTicket(t *testing.T) {
	ticket.HoldTicket()

}
func theTicketStatusShouldBe(arg1 string) error {

	return nil
}

func theUserShouldGetCkeckoutUrl() error {
	return nil
}

func ticketNumberOfBusNumberIs(arg1, arg2 int, arg3 string) error {

	return nil
}

func userRequestsReservation() error {
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Step(`^the ticket status should be "([^"]*)"$`, theTicketStatusShouldBe)
	ctx.Step(`^the user should get ckeckout url$`, theUserShouldGetCkeckoutUrl)
	ctx.Step(`^ticket number (\d+) of bus number (\d+) is "([^"]*)"$`, ticketNumberOfBusNumberIs)
	ctx.Step(`^user requests reservation$`, userRequestsReservation)
}
