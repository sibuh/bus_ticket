package routing

import (
	"event_ticket/internal/handler"

	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(group *gin.RouterGroup, handler handler.Ticket) {
	routes := []Route{
		{
			Method:  http.MethodPost,
			Path:    "/buy",
			Handler: handler.Buy,
		},
		{
			Method:  http.MethodPost,
			Path:    "/notify",
			Handler: handler.Notify,
		},
		{
			Method:  http.MethodGet,
			Path:    "/err",
			Handler: handler.Error,
		},
		{
			Method:  http.MethodGet,
			Path:    "/success",
			Handler: handler.Success,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/:nonce",
			Handler: handler.GetTicket,
		},
	}
	RegisterRoutes(group, routes)

}
