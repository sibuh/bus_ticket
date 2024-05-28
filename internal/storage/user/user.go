package user

import (
	"context"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"event_ticket/internal/storage"

	"golang.org/x/exp/slog"
)

type user struct {
	log     slog.Logger
	queries *db.Queries
}

func Init(logger slog.Logger, queries *db.Queries) storage.User {

	return &user{
		log:     logger,
		queries: queries,
	}
}
func (t *user) CreateUser(ctx context.Context, usr model.CreateUserRequest) (model.User, error) {
	u, err := t.queries.CreateUser(context.Background(), db.CreateUserParams{
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Phone:     usr.Phone,
		Email:     usr.Email,
		Username:  usr.Username,
		Password:  usr.Password,
	})
	if err != nil {
		t.log.Error("failed to register user", err)
		return model.User{}, err
	}
	return model.User{
		ID:        int(u.ID),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt.Time,
	}, nil

}

func (t *user) GetUser(ctx context.Context, username string) (model.User, error) {
	user, err := t.queries.GetUser(ctx, username)
	if err != nil {
		t.log.Error("failed to get user with the given id", err)
		return model.User{}, err
	}
	return model.User{
		ID:        int(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil

}
