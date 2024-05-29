package event

import (
	"context"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"event_ticket/internal/storage"
	"net/http"

	"golang.org/x/exp/slog"
)

type event struct {
	logger *slog.Logger
	db     *db.Queries
}

func Init(logger *slog.Logger, db *db.Queries) storage.Event {
	return &event{
		logger: logger,
		db:     db,
	}
}
func (e *event) PostEvent(ctx context.Context, eventParam model.Event) (model.Event, error) {
	ev, err := e.db.AddEvent(ctx, db.AddEventParams{
		Title:       eventParam.Title,
		Description: eventParam.Description,
		UserID:      int32(eventParam.UserID),
	})
	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to post event",
			RootError: err,
		}
		return model.Event{}, &newError
	}
	return model.Event{
		ID:    int(ev.ID),
		Title: ev.Title,
		Description: ev.Description,
		StartDate: ev.StartDate,
		EndDate: ev.EndDate,
		CreateAt: ev.CreatedAt,
		UpdatedAt: ev.UpdatedAt,
		DeletedAt: ev.DeletedAt.Time,
	}, nil
}
