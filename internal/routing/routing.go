package routing

import "github.com/gin-gonic/gin"

type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

func RegisterRoutes(group *gin.RouterGroup, routes []Route) {
	for _, r := range routes {
		group.Handle(r.Method, r.Path, r.Handler)
	}
}
