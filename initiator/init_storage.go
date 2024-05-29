package initiator

import (
	"event_ticket/internal/storage"
)

type Storage struct {
	user  storage.User
	event storage.Event
}

func NewStorage(user storage.User, event storage.Event) *Storage {
	return &Storage{
		user:  user,
		event: event,
	}
}
