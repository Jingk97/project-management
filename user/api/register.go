package api

import (
	"github.com/Jingk97/project-management-user/api/login"
	"github.com/Jingk97/project-management-user/model"
	"github.com/Jingk97/project-management-user/router"
)

// 用来上报所有的login路由函数

func RegisterApis() {
	handlerLogin := login.NewLoginHandler(model.RedisClient)
	getCaptchaRoute := &router.ModelRouteGroup{
		Path:    "/login/getCaptcha",
		Method:  "POST",
		Handler: handlerLogin.GetCaptcha,
	}
	getCaptchaRoute.AddRoute()
}
