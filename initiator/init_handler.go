package initiator

import (
	"event_ticket/internal/handler"
)

type Handler struct {
	user    handler.User
	payment handler.Payment
	event   handler.Event
}

func InitHandler(u handler.User, p handler.Payment, e handler.Event) *Handler {
	return &Handler{
		user:    u,
		payment: p,
		event:   e,
	}
}
