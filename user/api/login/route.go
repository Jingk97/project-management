package login

import (
	"github.com/Jingk97/project-management-user/router"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	log.Println("login路由已经添加到注册路由表中")
	router.AddRoute(&RouteUser{
		handler: NewHandlerLogin(),
	})
}

type RouteUser struct {
	handler *HandlerLogin
}

// Register RouteUser结构体在init就已经实体化
func (r *RouteUser) Register(engine *gin.RouterGroup) {
	userGroup := engine.Group("/login")
	{
		userGroup.POST("/getCaptcha", r.handler.getCaptcha)
	}
}
