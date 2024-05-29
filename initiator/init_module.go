package initiator

import "event_ticket/internal/module"

type Module struct {
	user  module.User
	event module.Event
}

func NewModule(user module.User, event module.Event) *Module {
	return &Module{
		user:  user,
		event: event,
	}
}
