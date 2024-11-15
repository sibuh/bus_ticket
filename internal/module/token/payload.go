package token

import (
	"bus_ticket/internal/data/db"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	TicketID uuid.UUID
	UserID   uuid.UUID
	ExpireAt time.Time
}

func NewTicketTokenPayload(userID, ticketID uuid.UUID, duration time.Duration) *Payload {
	return &Payload{
		TicketID: ticketID,
		UserID:   userID,
		ExpireAt: time.Now().Add(duration),
	}
}
func (pl *Payload) IsValid() error {
	if time.Now().Before(pl.ExpireAt) {
		return errors.New("token Expired")
	}
	return nil
}

type TicketInfo struct {
	db.User
	db.Ticket
}
