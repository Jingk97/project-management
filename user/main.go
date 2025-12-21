package main

import (
	"flag"
	"fmt"
	common "github.com/Jingk97/project-management-common"
	_ "github.com/Jingk97/project-management-user/api"
	"github.com/Jingk97/project-management-user/config"
	"github.com/Jingk97/project-management-user/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	var configFile string
	flag.StringVar(&configFile, "config-file", "", "指定启动加载配置文件路径")
	flag.Parse()
	// 初始化所有中间件以及切面服务
	cfg := initServer(configFile)
	r.Use(common.GinLogger(), common.GinRecovery(true))
	// 初始化注册路由
	router.InitRouter(r)
	common.Run(r, fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port), cfg.ServerName)
}

func initServer(configFile string) *config.Config {
	// 加载配置文件；最高优先级是启动命令时候获取的参数
	// TODO：后期使用启动参数config-file，如果值file则是本地配置文件，如果是nacos则是读取nacos配置中心文件
	// 初始化加载配置内容
	cfg := config.InitConfig(configFile)
	// 初始化加载zap日志库
	cfg.InitZapLog()
	// 初始化加载redis
	cfg.InitRedisDB()
	return cfg
}
