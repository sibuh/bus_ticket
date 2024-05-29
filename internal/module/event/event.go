package event

import (
	"context"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"event_ticket/internal/storage"

	"golang.org/x/exp/slog"
)

type event struct {
	logger *slog.Logger
	es     storage.Event
}

func Init(logger *slog.Logger, es storage.Event) module.Event {
	return &event{
		logger: logger,
		es:     es,
	}
}
func (e *event) PostEvent(ctx context.Context, postEvent model.Event) (model.Event, error) {

	return e.es.PostEvent(ctx, postEvent)

}

func (e *event) FetchEvents(ctx context.Context) ([]model.Event, error) {

	return e.es.FetchEvents(ctx)

}
