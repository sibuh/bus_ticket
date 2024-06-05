package initiator

import (
	"event_ticket/internal/handler"
)

type Handler struct {
	user    handler.User
	payment handler.Payment
	event   handler.Event
	ticket  handler.Ticket
}

func InitHandler(u handler.User, p handler.Payment, e handler.Event, t handler.Ticket) *Handler {
	return &Handler{
		user:    u,
		payment: p,
		event:   e,
		ticket:  t,
	}
}
