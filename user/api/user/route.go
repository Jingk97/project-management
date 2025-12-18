package user

import (
	"github.com/Jingk97/project-management-user/router"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	log.Println("init user router")
	router.Register(&RouteUser{})
}

type RouteUser struct{}

func (*RouteUser) Route(engine *gin.Engine) {
	h := HandlerUser{}
	engine.POST("/project/login/getCaptcha", h.getCaptcha)
}
