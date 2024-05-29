package event

import (
	"context"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"event_ticket/internal/storage"
	"net/http"
	"strings"

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
		UserID:      eventParam.UserID,
		StartDate:   eventParam.StartDate,
		EndDate:     eventParam.EndDate,
		Price:       eventParam.Price,
	})
	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusInternalServerError,
			Message:   "failed to post event",
			RootError: err,
		}
		e.logger.Error(newError.Error(), newError)
		return model.Event{}, &newError
	}
	return model.Event{
		ID:          ev.ID,
		Title:       ev.Title,
		Description: ev.Description,
		UserID:      ev.UserID,
		StartDate:   ev.StartDate,
		EndDate:     ev.EndDate,
		Price:       ev.Price,
		CreateAt:    ev.CreatedAt,
		UpdatedAt:   ev.UpdatedAt,
		DeletedAt:   ev.DeletedAt.Time,
	}, nil
}

func (e *event) FetchEvents(ctx context.Context) ([]model.Event, error) {
	evs, err := e.db.FetchEvents(ctx)
	if err != nil {
		var message, code = func() (string, int) {
			if strings.Contains(err.Error(), "not found") {
				return "resource not found", http.StatusNotFound
			} else {
				return "unable to get events", http.StatusInternalServerError
			}
		}()
		newError := model.Error{
			ErrCode:   code,
			Message:   message,
			RootError: err,
		}
		return nil, &newError
	}

	var events []model.Event
	for _, ev := range evs {
		events = append(events, model.Event{
			ID:          ev.ID,
			Title:       ev.Title,
			Description: ev.Description,
			StartDate:   ev.StartDate,
			EndDate:     ev.EndDate,
			Price:       ev.Price,
			CreateAt:    ev.CreatedAt,
			UpdatedAt:   ev.UpdatedAt,
			DeletedAt:   ev.DeletedAt.Time,
		})
	}
	return events, nil
}
