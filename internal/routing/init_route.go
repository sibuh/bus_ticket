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
			Path:    "/create",
			Handler: u.CreateUser,
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
