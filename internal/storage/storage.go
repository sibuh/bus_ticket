package storage

import (
	"context"
	"event_ticket/internal/model"
)

type User interface {
	CreateUser(ctx context.Context, usr model.CreateUserRequest) (model.User, error)
	GetUser(ctx context.Context, id int32) (model.User, error)
}
