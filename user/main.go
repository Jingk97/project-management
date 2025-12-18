package main

import (
	common "github.com/Jingk97/project-management-common"
	_ "github.com/Jingk97/project-management-user/api"
	"github.com/Jingk97/project-management-user/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 初始化注册路由
	router.InitRouter(r)
	common.Run(r, "127.0.0.1:8080", "project-user")
}
