package itinerary

import (
	"bus_ticket/internal/data/db"
	"bus_ticket/internal/module"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// You should only set `manual=true` for trips created by sys admin
// if set to true, it will not check existing line data to create trips
func (l *line) CreateTrip(ctx context.Context, tripPayload db.CreateTripParams, manual bool) (db.LineTrip, error) {
	if manual == true {
		result, err := l.q.CreateTrip(ctx, tripPayload)
		if err != nil {
			return db.LineTrip{}, err
		}
		return result, nil
	}

	valid, err := l.IsTimeValidForLine(ctx, tripPayload.Line, tripPayload.Date)
	if err != nil {
		return db.LineTrip{}, err
	}

	if !valid {
		return db.LineTrip{}, errors.New("Invalid date and time")
	}
	diff := tripPayload.Date.Sub(time.Now())
	// config := l.GetClientConfig(uuid.New())
	// max_res_day := config.max_ticket_reservation_day

	if !(diff.Hours() > time.Hour.Hours()*24) {
		return db.LineTrip{}, errors.New("Try in later days")
	}
	trip, err := l.q.CreateTrip(ctx, tripPayload)
	return trip, nil
}

// TODO: What happens if a line schedule or price is updated after a trip is created
// I believe line update should be applicable only after some day
func (l *line) UpdateTrip(ctx context.Context, lineID string) (db.LineTrip, error) {
	// assign bus and/or driver/redat or so.
	// fabricated
	return db.LineTrip{}, nil
}

// what if a trip is canceled by sys admin

type ClientConfig struct {
	max_ticket_reservation_day int
}

func (l *line) GetClientConfig(clientId uuid.UUID) *ClientConfig {
	// query the data base using client config
	return &ClientConfig{
		max_ticket_reservation_day: 10,
	}
}

func (l *line) IsTimeValidForLine(ctx context.Context, lineID uuid.UUID, t time.Time) (bool, error) {
	line, err := l.q.GetLine(ctx, lineID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, &module.NotFoundError{
				ResourceName: "line",
				ResourceId:   lineID.String(),
				Context:      "creating trip",
			}
		}
		return false, &module.UnknownServerError{
			Message: "Uknown server error happened while creating trip",
			Cause:   err,
		}
	}

	jsonBlob, err := line.Schedule.MarshalJSON()
	if err != nil {
		return false, &module.UnknownServerError{Message: "Couldn't read schedule from db", Cause: err}
	}

	var schedule Schedule
	unmarshalErr := json.Unmarshal(jsonBlob, &schedule)

	if unmarshalErr != nil {
		return false, errors.New("Unmarshal Error ")
	}

	_, _, day := t.Date()

	times := schedule[time.Weekday(day)]

	for _, val := range times {
		if val.Equal(t) {
			return true, nil
		}
	}
	return false, nil
}
