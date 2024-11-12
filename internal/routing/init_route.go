package routing

import (
	"bus_ticket/internal/handler"
	"bus_ticket/internal/middleware"

	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(group *gin.RouterGroup, u handler.User, t handler.Ticket, md middleware.Middleware) {
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
			Path:    "/token",
			Handler: u.RefreshToken,
			Mwares:  []gin.HandlerFunc{md.Authenticate()},
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/:intent_id",
			Handler: t.GetTicket,
			Mwares:  []gin.HandlerFunc{md.Authenticate()},
		},
	}
	RegisterRoutes(group, routes)

}
