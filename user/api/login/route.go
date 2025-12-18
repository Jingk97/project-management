package user

import (
	"github.com/gin-gonic/gin"
)

type RouteUser struct {
	handler *HandlerUser
}

func (r *RouteUser) Route(engine *gin.RouterGroup) {
	userGroup := engine.Group("/user")
	{
		userGroup.POST("/login/getCaptcha", r.handler.getCaptcha)

	}
}
