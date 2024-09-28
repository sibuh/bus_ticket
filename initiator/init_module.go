package initiator

import "event_ticket/internal/module"

type Module struct {
	user    module.User
	event   module.Event
	payment module.Payment
	ticket  module.Ticket
}

func NewModule(user module.User, payment module.Payment, tkt module.Ticket) *Module {
	return &Module{
		user:    user,
		payment: payment,
		ticket:  tkt,
	}
}
