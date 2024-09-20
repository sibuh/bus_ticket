package callback

import (
	"event_ticket/internal/model"
	"event_ticket/internal/module/scheduler"
)

type Callback struct {
	scheduler scheduler.Scheduler
}

func Init(scheduler scheduler.Scheduler) *Callback {
	return &Callback{
		scheduler: scheduler,
	}
}

func (c *Callback) exitScheduler(payload model.Payment) {
	sessionId := payload.IntentID

	ch := c.scheduler.Get(sessionId)
	ch <- sessionId
}

func (c *Callback) HandlePaymentStatusUpdate(payload model.Payment) {
	c.exitScheduler(payload)
}
