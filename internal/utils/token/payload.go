package token

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID       uuid.UUID
	Username string
	IssuedAt time.Time
	ExpireAt time.Time
}

func NewPayload(username string, duration time.Duration) *Payload {
	return &Payload{
		ID:       uuid.New(),
		Username: username,
		IssuedAt: time.Now(),
		ExpireAt: time.Now().Add(duration),
	}
}
func (pl *Payload) Valid() bool {
	return time.Now().Before(pl.ExpireAt)
}
