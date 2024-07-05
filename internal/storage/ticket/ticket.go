package ticket

import (
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"event_ticket/internal/storage"

	"golang.org/x/exp/slog"
)

type ticket struct {
	logger *slog.Logger
	*db.Queries
}

func Init(logger *slog.Logger) storage.Ticket {
	return &ticket{
		logger: logger,
	}
}
func (t *ticket) HoldTicket(ticketNo, tripId int32) (model.Ticket, error) {

	return model.Ticket{}, nil
}
func (t *ticket) GetTicket(tktNo, tripId int32) (model.Ticket, error) {
	return model.Ticket{}, nil
}
func (t *ticket) UnholdTicket(tktNo, tripID int32) (model.Ticket, error) {
	return model.Ticket{}, nil
}
