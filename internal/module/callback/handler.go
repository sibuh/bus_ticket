package callback

import (
	"event_ticket/internal/model"
	"event_ticket/internal/module/scheduler"
)

type Callback struct {
	schedulerMap scheduler.Scheduler
}

func (c *Callback) exitScheduler(payload model.Payment) {
	sessionId := payload.IntentID

	ch := c.schedulerMap.Map[sessionId]
	ch <- sessionId

	c.schedulerMap.Remove(sessionId)
}

func (c *Callback) handlePaymentStatusUpdate(payload model.Payment) {
	c.exitScheduler(payload)

	// do callback business logic
}
