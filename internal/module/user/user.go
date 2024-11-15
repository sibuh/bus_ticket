package user

import (
	"bus_ticket/internal/data/db"
	"bus_ticket/internal/model"
	"bus_ticket/internal/module"
	"bus_ticket/internal/utils/token"
	"bus_ticket/pkg/pass"
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

type user struct {
	logger        *slog.Logger
	tokenMaker    token.TokenMaker
	q             db.Querier
	tokenDuration time.Duration
}

func Init(logger *slog.Logger, q db.Querier, tokenMaker token.TokenMaker, td time.Duration) module.User {

	return &user{
		logger:        logger,
		tokenMaker:    tokenMaker,
		q:             q,
		tokenDuration: td,
	}

}

func (u *user) CreateUser(ctx context.Context, usr model.CreateUserRequest) (db.User, error) {
	if err := usr.Validate(); err != nil {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "invalid input",
			RootError: err,
		}
		u.logger.Info("invalid user input", newError)
		return db.User{}, &newError
	}
	hash, err := pass.HashPassword(usr.Password)
	if err != nil {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "failed to hash password",
			RootError: err,
		}
		u.logger.Error("failed to hash password", newError)
		return db.User{}, &newError
	}
	usr.Password = hash

	return u.q.CreateUser(ctx, db.CreateUserParams{
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Username:  usr.Username,
		Phone:     usr.Phone,
		Password:  usr.Password,
		Email:     usr.Email,
	})
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

	usr, err := u.q.GetUser(ctx, logReq.Username)
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
	authTokenPayload := token.NewAuthTokenPayload(logReq.Username, u.tokenDuration)
	return u.tokenMaker.CreateToken(authTokenPayload)
}

func (u *user) RefreshToken(ctx context.Context, username string) (string, error) {
	authTokenPayload := token.NewAuthTokenPayload(username, u.tokenDuration)
	return u.tokenMaker.CreateToken(authTokenPayload)
}
