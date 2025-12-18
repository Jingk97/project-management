package login

import "github.com/gin-gonic/gin"

type HandlerLogin struct {
}

func (*HandlerLogin) getCaptcha(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "getCaptcha success",
	})
}
