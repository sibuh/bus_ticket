package initiator

import (
	"event_ticket/internal/handler"
)

type Handler struct {
	user    handler.User
	payment handler.Payment
}

func InitHandler(u handler.User, p handler.Payment) *Handler {
	return &Handler{
		user:    u,
		payment: p,
	}
}
