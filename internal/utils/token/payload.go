package token

import (
	"errors"
	"time"
)

type Payload struct {
	UserID   string
	IssuedAt time.Time
	ExpireAt time.Time
}

func NewAuthTokenPayload(userId string, duration time.Duration) TokenValidator {
	return &Payload{
		UserID:   userId,
		IssuedAt: time.Now(),
		ExpireAt: time.Now().Add(duration),
	}
}
func (pl *Payload) IsValid() error {
	if !time.Now().Before(pl.ExpireAt) {
		return errors.New("token expired")
	}
	return nil
}
