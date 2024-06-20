package ticket

import (
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"event_ticket/internal/storage"

	"golang.org/x/exp/slog"
)

type ticket struct {
	log           slog.Logger
	storageTicket storage.Ticket
}

func Init(log slog.Logger, tkt storage.Ticket) module.Ticket {
	return &ticket{
		log:           log,
		storageTicket: tkt,
	}
}
func (t *ticket) HoldTicket(tktNo, tripId int32) (model.Ticket, error) {
	tkt, err := t.storageTicket.HoldTicket(tktNo, tripId)
	if err != nil {
		return model.Ticket{}, err
	}
	return tkt, nil
}
