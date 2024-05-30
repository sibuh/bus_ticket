package storage

import (
	"context"
	"event_ticket/internal/model"
)

type User interface {
	CreateUser(ctx context.Context, usr model.CreateUserRequest) (model.User, error)
	GetUser(ctx context.Context, username string) (model.User, error)
}

type Event interface {
	PostEvent(ctx context.Context, eventParam model.Event) (model.Event, error)
	FetchEvents(ctx context.Context) ([]model.Event, error)
	FetchEvent(ctx context.Context, id int32) (model.Event, error)
}
