package callback

import (
	"bus_ticket/internal/model"
	"bus_ticket/internal/module/schedule"
)

type Callback struct {
	Scheduler *schedule.Scheduler
}

func Init(scheduler *schedule.Scheduler) *Callback {
	return &Callback{
		Scheduler: scheduler,
	}
}

func (c *Callback) exitScheduler(payload model.Payment) {
	sessionId := payload.IntentID

	ch := c.Scheduler.Get(sessionId)
	if ch != nil {
		ch <- sessionId
	}

}

func (c *Callback) HandlePaymentStatusUpdate(payload model.Payment) {
	c.exitScheduler(payload)
	//TODO:do databse update
}
