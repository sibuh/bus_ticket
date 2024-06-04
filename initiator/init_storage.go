package initiator

import (
	"event_ticket/internal/storage"
)

type Storage struct {
	user  storage.User
	event storage.Event
	pmt   storage.Payment
}

func NewStorage(user storage.User, event storage.Event, pmt storage.Payment) *Storage {
	return &Storage{
		user:  user,
		event: event,
		pmt:   pmt,
	}
}
