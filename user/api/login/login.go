package user

import "github.com/gin-gonic/gin"

type HandlerUser struct {
}

func (*HandlerUser) getCaptcha(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "getCaptcha success",
	})
}
