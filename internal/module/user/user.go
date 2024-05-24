package user

import (
	"context"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"event_ticket/internal/storage"
	"event_ticket/internal/utils/pass"
	"event_ticket/internal/utils/token"
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"
)

type user struct {
	logger     slog.Logger
	user       storage.User
	tokenMaker token.TokenMaker
}

func Init(logger slog.Logger, usr storage.User, tokenMaker token.TokenMaker) module.User {

	return &user{
		logger:     logger,
		user:       usr,
		tokenMaker: tokenMaker,
	}

}

func (u *user) CreateUser(ctx context.Context, usr model.CreateUserRequest) (model.User, error) {
	if err := usr.Validate(); err != nil {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "invalid input",
			RootError: err,
		}
		u.logger.Info("invalid user input", newError)
		return model.User{}, &newError
	}
	hash, err := pass.HashPassword(usr.Password)
	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "failed to hash password",
			RootError: err,
		}
		u.logger.Error("failed to hash password", newError)
		return model.User{}, &newError
	}
	usr.Password = hash

	return u.user.CreateUser(ctx, usr)
}

func (u *user) LoginUser(ctx context.Context, logReq model.LoginRequest) (string, error) {
	if err := logReq.Validate(); err != nil {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "invalid user input",
			RootError: err,
		}
		u.logger.Info("invalid user input", newError.RootError.Error())
		return "", &newError
	}

	usr, err := u.user.GetUser(ctx, logReq.Username)
	if err != nil {
		return "", err
	}
	if !pass.CheckHashPassword(logReq.Password, usr.Password) {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "incorrect password",
			RootError: fmt.Errorf("invalid input"),
		}
		return "", &newError
	}

	return u.tokenMaker.CreateToken(usr.Username)
}
