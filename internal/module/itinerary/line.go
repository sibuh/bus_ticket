package itinerary

import (
	"context"
	"database/sql"
	"errors"
	"event_ticket/internal/data/db"
	"event_ticket/internal/module"
	"fmt"
	"time"

	"golang.org/x/exp/slog"
)

type line struct {
	logger slog.Logger
	q      db.Querier
}

type Schedule map[time.Weekday][]time.Time

// {1: [10, 5], 2: [5]}

func (l *line) CreateLine(ctx context.Context, paylaod db.CreateLineParams) (db.Line, error) {
	// Departure and Destination should be enum fields
	result, err := l.q.CreateLine(ctx, paylaod)
	if err != nil {
		return db.Line{}, err
	}
	return result, nil
}

// Assumed only price and schedule would be updated,
// Should create new line instead of updating destination and departure place
func (l *line) UpdateLine(ctx context.Context, payload db.UpdateLineParams) (db.Line, error) {
	result, err := l.q.UpdateLine(ctx, payload)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Line{}, &module.NotFoundError{
				ResourceName: "line",
				ResourceId:   payload.ID.String(),
				Context:      "updating line",
			}
		}
		l.logger.Error(fmt.Sprintf("query couldn't update line %s", payload.ID), payload)
		return db.Line{}, &module.UnknownServerError{
			Message: fmt.Sprintf("Couldn't update line %s", payload.ID),
			Cause:   err,
		}
	}

	return result, nil
}
