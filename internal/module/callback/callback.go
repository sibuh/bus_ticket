package callback

import (
	"event_ticket/internal/model"
	"event_ticket/internal/module/schedule"
)

type Callback struct {
	Scheduler schedule.Scheduler
}

func Init() *Callback {
	return &Callback{
		Scheduler: *schedule.Init(),
	}
}

func (c *Callback) ExitScheduler(payload model.Payment) {
	sessionId := payload.IntentID

	ch := c.Scheduler.Get(sessionId)
	ch <- sessionId
}

func (c *Callback) HandlePaymentStatusUpdate(payload model.Payment) {
	c.ExitScheduler(payload)
	//TODO:do databse update
}
