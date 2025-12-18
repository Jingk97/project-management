package router

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Register(engin *gin.RouterGroup)
}

var routes []Router

func AddRoute(r Router) {
	routes = append(routes, r)
}

func InitRouter(ginEngine *gin.Engine) {
	apiv1 := ginEngine.Group("/api/v1")
	for _, route := range routes {
		route.Register(apiv1)
	}
}
