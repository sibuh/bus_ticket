package ticket

import (
	"event_ticket/internal/model"
	"fmt"
)

type MockStorageTicket struct {
	Tkt model.Ticket
}

func InitMock(tkt model.Ticket) *MockStorageTicket {
	return &MockStorageTicket{Tkt: tkt}
}
func (m *MockStorageTicket) HoldTicket(ticketNo, tripId int32) (model.Ticket, error) {

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
