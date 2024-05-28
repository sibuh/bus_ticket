package user

import (
	"event_ticket/internal/handler"
	"event_ticket/internal/model"
	"event_ticket/internal/module"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type user struct {
	logger slog.Logger
	user   module.User
}

func Init(logger slog.Logger, usr module.User) handler.User {
	return &user{
		logger: logger,
		user:   usr,
	}
}

func (u *user) CreateUser(c *gin.Context) {
	var user model.CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "failed to marshal response body",
			RootError: err,
		}
		c.JSON(http.StatusBadRequest, newError)
		return
	}
	createdUser, err := u.user.CreateUser(c, user)

	if err != nil {
		newError := err.(*model.Error)
		c.JSON(newError.ErrCode, err)
		return
	}
	c.JSON(http.StatusOK, createdUser)

}
func (u *user) LoginUser(c *gin.Context) {
	var logReq model.LoginRequest
	if err := c.ShouldBind(&logReq); err != nil {

		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "failed to bind the request body",
			RootError: err,
		}
		u.logger.Error("failed to bind the request body", newError)
		c.JSON(newError.ErrCode, newError)
		return
	}
	token, err := u.user.LoginUser(c, logReq)
	if err != nil {
		newErr := err.(*model.Error)
		c.JSON(newErr.ErrCode, newErr)
		return
	}
	c.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: token})
}
