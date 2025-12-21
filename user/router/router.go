package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// RouteGroup
/*
路由设计：/api/v1 api+version进行拼接后续可以进行版本迭代
后续/api/v1/业务模块/具体函数进行拼接；例如 /api/v1/login/getCaptcha
因此RouteGroup 对应业务模块
Group：业务模块名称；Path是模块后接口函数；Method是后续接口方法；Handler就是对应方法
*/
var allRoutes []*ModelRouteGroup

type ModelRouteGroup struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

func (r *ModelRouteGroup) AddRoute() {
	allRoutes = append(allRoutes, r)
}

func InitRouters(r *gin.Engine) {
	apiv1 := r.Group("/api/v1")
	fmt.Println("register all routes", allRoutes[0])
	for _, route := range allRoutes {
		switch route.Method {
		case "GET":
			apiv1.GET(route.Path, route.Handler)
		case "POST":
			apiv1.POST(route.Path, route.Handler)
		case "PUT":
			apiv1.PUT(route.Path, route.Handler)
		case "DELETE":
			apiv1.DELETE(route.Path, route.Handler)
		default:
			apiv1.GET(route.Path, route.Handler)
		}
	}
}
