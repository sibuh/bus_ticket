package user

import (
	"bus_ticket/internal/data/db"
	"bus_ticket/internal/handler"
	"bus_ticket/internal/model"
	"bus_ticket/internal/module"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type user struct {
	logger *slog.Logger
	user   module.User
}

func Init(logger *slog.Logger, usr module.User) handler.User {
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
	c.Header("Authorization", token)
}

func (u *user) RefreshToken(c *gin.Context) {

	user := c.Value("user").(db.User)
	if user.Username == "" {
		newError := model.Error{
			ErrCode:   http.StatusBadRequest,
			Message:   "username not set to context",
			RootError: fmt.Errorf("username not set to context"),
		}
		c.JSON(newError.ErrCode, newError)
		return
	}
	token, err := u.user.RefreshToken(c, user.Username)
	if err != nil {
		newError := err.(*model.Error)
		c.JSON(newError.ErrCode, newError.Message)
		return
	}
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, nil)
}
