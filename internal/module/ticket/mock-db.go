package ticket

import (
	"context"
	"database/sql"
	"event_ticket/internal/data/db"
)

type MockQueries struct {
	db.Querier
	Tkt db.Ticket
}

func (m MockQueries) UpdateTicketStatus(ctx context.Context, arg db.UpdateTicketStatusParams) (db.Ticket, error) {
	m.Tkt = db.Ticket{
		TicketNo: arg.TicketNo,
		BusNo:    arg.BusNo,
		TripID:   arg.TripID,
		Status: sql.NullString{
			String: "Onhold",
			Valid:  true,
		}}

	return m.Tkt, nil
}
