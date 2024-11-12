package initiator

import "bus_ticket/internal/module"

type Module struct {
	user   module.User
	ticket module.Ticket
}

func NewModule(user module.User, tkt module.Ticket) *Module {
	return &Module{
		user:   user,
		ticket: tkt,
	}
}
