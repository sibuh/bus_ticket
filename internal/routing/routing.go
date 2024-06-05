package routing

import "github.com/gin-gonic/gin"

type Route struct {
	Method  string
	Path    string
	Mwares  []gin.HandlerFunc
	Handler gin.HandlerFunc
}

func RegisterRoutes(group *gin.RouterGroup, routes []Route) {
	for _, r := range routes {
		handlers := []gin.HandlerFunc{}
		if len(r.Mwares) > 0 {
			handlers = append(handlers, r.Mwares...)
		}
		handlers = append(handlers, r.Handler)
		group.Handle(r.Method, r.Path, handlers...)
	}
}
