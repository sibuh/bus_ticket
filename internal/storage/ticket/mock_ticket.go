package ticket

import (
	"event_ticket/internal/model"
)

type MockStorageTicket struct {
	Tkt model.Ticket
}

func InitMock(tkt model.Ticket) *MockStorageTicket {
	return &MockStorageTicket{Tkt: tkt}
}
func (m *MockStorageTicket) HoldTicket(ticketNo, tripId int32) (model.Ticket, error) {
	m.Tkt.Status = "onhold"
	return m.Tkt, nil

}
func (m *MockStorageTicket) AddTicket(ticketNo, busNo, tripId int32, status string) (model.Ticket, error) {
	m.Tkt = model.Ticket{
		TripId:   tripId,
		TicketNo: ticketNo,
		BusNo:    busNo,
		Status:   status,
	}
	return m.Tkt, nil
}

func (m *MockStorageTicket) CheckTicketStatus(tktNo int32) string {
	if tktNo == m.Tkt.TicketNo {
		return m.Tkt.Status
	}
	return ""
}
