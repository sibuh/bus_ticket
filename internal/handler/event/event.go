package event

import (
	"event_ticket/internal/handler"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type event struct {
	logger *slog.Logger
	em     module.Event
}

func Init(logger *slog.Logger, em module.Event) handler.Event {
	return &event{
		logger: logger,
		em:     em,
	}
}
func (e *event) PostEvent(c *gin.Context) {

	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "failed to bind request body",
			RootError: err,
		}
		e.logger.Info("failed to marshal request body", err.Error())
		c.JSON(newError.ErrCode, newError)
		return
	}
	ev, err := e.em.PostEvent(c, event)
	if err != nil {
		e := err.(*model.Error)
		c.JSON(e.ErrCode, e)
		return
	}
	c.JSON(http.StatusOK, ev)
}

func (e *event) FetchEvents(c *gin.Context) {

	events, err := e.em.FetchEvents(c)
	if err != nil {
		e := err.(*model.Error)
		c.JSON(e.ErrCode, e)
		return
	}
	//TODO: pagination for FetchEvents
	c.JSON(http.StatusOK, events)
}
