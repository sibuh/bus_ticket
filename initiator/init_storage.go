package initiator

import (
	"event_ticket/internal/storage"
)

type Storage struct {
	user storage.User
}

func NewStorage(user storage.User) *Storage {
	return &Storage{
		user: user,
	}
}
