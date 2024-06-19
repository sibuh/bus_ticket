package ticket

import (
	"event_ticket/internal/module"

	"golang.org/x/exp/slog"
)

type ticket struct {
	log slog.Logger
}

func Init(log slog.Logger) module.Ticket {
	return &ticket{
		log: log,
	}
}
func HoldTicket() {}
