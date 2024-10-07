package initiator

import (
	"event_ticket/internal/handler"
)

type Handler struct {
	user   handler.User
	ticket handler.Ticket
}

func InitHandler(u handler.User, t handler.Ticket) *Handler {
	return &Handler{
		user:   u,
		ticket: t,
	}
}
