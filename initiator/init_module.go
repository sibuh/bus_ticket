package initiator

import "event_ticket/internal/module"

type Module struct {
	user    module.User
	event   module.Event
	payment module.Payment
}

func NewModule(user module.User, event module.Event, payment module.Payment) *Module {
	return &Module{
		user:    user,
		event:   event,
		payment: payment,
	}
}
