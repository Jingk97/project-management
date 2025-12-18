package router

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Route(engin *gin.Engine)
}

type RegisterRouter struct {
}

//func NewRegisterRouter() *RegisterRouter {
//	return &RegisterRouter{}
//}

func (r *RegisterRouter) Route(router Router, engin *gin.Engine) {
	router.Route(engin)
}

var routes []Router

func Register(ro ...Router) {
	routes = append(routes, ro...)
}

func InitRouter(ginEngine *gin.Engine) {
	//r := NewRegisterRouter()
	//r.Route(&user.RouteUser{}, ginEngine)

	for _, route := range routes {
		route.Route(ginEngine)
	}
}
