package login

import (
	common "github.com/Jingk97/project-management-common"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type HandlerLogin struct {
}

func (*HandlerLogin) getCaptcha(ctx *gin.Context) {
	resp := common.Result{}
	// 获取请求参数（主要是手机号）
	mobileNum := ctx.PostForm("mobile")

	// 校验手机号是否符合规则
	if mobileNum == "" {
		ctx.JSON(http.StatusOK, resp.Fail(common.IllegalPhoneNumber, "Phone Number Miss!"))
		log.Printf("请求参数，手机号为空")
		return
	}
	if !common.IsValidateMobile(mobileNum) {
		ctx.JSON(http.StatusOK, resp.Fail(common.IllegalPhoneNumber, "Phone Number Illegal!"))
		log.Printf("请求参数，手机号：%s，无效", mobileNum)
		return
	}
	// 生成验证码（验证码需要协程发送）
	code, err := common.GenerateCode(6)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, resp.Fail(common.GenerateCodeWrong, err.Error()))
		return
	}
	// 发送验证码
	go func() {
		time.Sleep(400 * time.Millisecond)
		log.Println("调用短信发送成功；code：", code)
		ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
		defer cancel()
		err := h.cache.Put(ctx, "REGISTER_"+mobileNum, code, 15*time.Minute)
		if err != nil {
			log.Println("验证码入redis出错：", err)
			return
		}
		log.Println("已经将code入redis缓存", "REGISTER_"+mobileNum, ": ", code)
	}()
	// 存入redis
	// 返回结果
	ctx.JSON(200, resp.Success(code))
	return
}
