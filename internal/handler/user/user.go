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
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	createdUser, err := u.user.CreateUser(c, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, createdUser)

}
func (u *user) LoginUser(c *gin.Context) {

}
