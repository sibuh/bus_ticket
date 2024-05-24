package user

import (
	"context"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"event_ticket/internal/storage"
	"event_ticket/internal/utils/pass"
	"fmt"

	"golang.org/x/exp/slog"
)

type user struct {
	logger slog.Logger
	user   storage.User
}

func Init(logger slog.Logger, usr storage.User) module.User {
	return &user{
		logger: logger,
		user:   usr,
	}

}

func (u *user) CreateUser(ctx context.Context, usr model.CreateUserRequest) (model.User, error) {
	if err := usr.Validate(); err != nil {
		return model.User{}, err
	}
	hash, err := pass.HashPassword(usr.Password)
	if err != nil {
		return model.User{}, err
	}
	fmt.Println("request data:", usr, usr.Password)
	usr.Password = hash
	createdUser, err := u.user.CreateUser(ctx, usr)
	if err != nil {
		u.logger.Error("failed to create user", err)
		return model.User{}, err
	}
	return createdUser, nil

}
func (u *user) GetUser(ctx context.Context, id int32) (model.User, error) {

	return u.user.GetUser(ctx, id)
}
func (u *user) LoginUser(ctx context.Context) (string, error) {
	return "", nil

}
