package login

import (
	"context"
	common "github.com/Jingk97/project-management-common"
	"github.com/Jingk97/project-management-user/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type HandlerLogin struct {
	cache *model.RedisCache
}

// NewLoginHandler 还是要有这个new函数，结构体后续变成有状态后，上一层级需要将有状态的接口实体化
func NewLoginHandler(redis *model.RedisCache) *HandlerLogin {
	return &HandlerLogin{
		cache: redis,
	}
}

func (h *HandlerLogin) GetCaptcha(ctx *gin.Context) {
	resp := common.Result{}
	// 获取请求参数（主要是手机号）
	mobileNum := ctx.PostForm("mobile")

	// 校验手机号是否符合规则
	if mobileNum == "" {
		ctx.JSON(http.StatusOK, resp.Fail(common.IllegalPhoneNumber, "Phone Number Miss!"))
		zap.L().Warn("请求参数，手机号为空")
		return
	}
	if !common.IsValidateMobile(mobileNum) {
		ctx.JSON(http.StatusOK, resp.Fail(common.IllegalPhoneNumber, "Phone Number Illegal!"))
		zap.L().Warn("用户携带手机号不合法：", zap.String("mobile", mobileNum))
		return
	}
	// 生成验证码（验证码需要协程发送）
	code, err := common.GenerateCode(6)
	if err != nil {
		zap.L().Warn("验证码生成失败", zap.String("mobile", mobileNum), zap.Error(err))
		ctx.JSON(http.StatusOK, resp.Fail(common.GenerateCodeWrong, err.Error()))
		return
	}
	// 发送验证码
	go func() {
		time.Sleep(400 * time.Millisecond)
		zap.L().Info("调用短信发送成功；code：", zap.String("code", code))
		ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
		defer cancel()
		// 存入redis
		err = h.cache.Put(ctx, "REGISTER_"+mobileNum, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("验证码落入redis出错：", zap.String("mobile", mobileNum), zap.Error(err))
			return
		}
		zap.L().Info("已经将code入redis缓存,REGISTER_", zap.String("mobile", mobileNum), zap.String("code", code))
	}()
	// 返回结果
	ctx.JSON(200, resp.Success(code))
	return
}
