package storage

import (
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
)

type Ticket interface {
	RegisterUserToDb(user model.User, sessionid, nonce string) error
	UpdatePaymentStatus(status, sid string) (db.User, error)
	GetUser(nonce string) (db.User, error)
}
