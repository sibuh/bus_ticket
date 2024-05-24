package initiator

import "event_ticket/internal/module"

type Module struct {
	user module.User
}

func NewModule(user module.User) *Module {
	return &Module{
		user: user,
	}
}
