package ticket

import (
	"event_ticket/internal/model"
	"event_ticket/internal/storage"

	"golang.org/x/exp/slog"
)

type ticket struct {
	logger slog.Logger
}

func Init(logger slog.Logger) storage.Ticket {
	return &ticket{
		logger: logger,
	}
}
func (t *ticket) ReserveTicket(ticketNo, tripId int32) (model.Ticket, error) {
	return model.Ticket{}, nil
}
func (t *ticket) GetTicket(tktNo, tripId int32) (model.Ticket, error) {
	return model.Ticket{}, nil
}
func (t *ticket) UnholdTicket(tktNo, tripID int32) (model.Ticket, error) {
	return model.Ticket{}, nil
}
