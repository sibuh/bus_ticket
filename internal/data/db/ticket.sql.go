package db

import (
	"context"
)

type UpdateTicketStatusParams struct {
	TicketNo int32
	BusNo    int32
	TripID   int32
}

func (q *Queries) UpdateTicketStatus(ctx context.Context, arg UpdateTicketStatusParams) (Ticket, error) {
	return Ticket{}, nil
}
