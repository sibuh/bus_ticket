package routing

import (
	"event_ticket/internal/handler"

	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(group *gin.RouterGroup, u handler.User, p handler.Payment) {
	routes := []Route{
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: u.CreateUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: u.LoginUser,
		},

		{
			Method:  http.MethodGet,
			Path:    "/pk",
			Handler: p.GetPublishableKey,
		},
		{
			Method:  http.MethodGet,
			Path:    "/cpi",
			Handler: p.HandleCreatePaymentIntent,
		},
	}
	RegisterRoutes(group, routes)

}
