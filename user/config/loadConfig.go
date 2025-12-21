package config

// 整体配置文件都从common包中copy过来，common包中代码为通案例

import (
	common "github.com/Jingk97/project-management-common"
	"github.com/Jingk97/project-management-user/model"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	ServerInfo `mapstructure:"server"`
	viper      *viper.Viper
}

type ServerInfo struct {
	ServerName string `mapstructure:"serverName"`
	Mode       string `mapstructure:"mode"`
	Addr       string `mapstructure:"addr"`
	Port       uint16 `mapstructure:"port"`
}

func InitConfig(configFile string) *Config {
	v := viper.New()
	conf := &Config{
		viper:      v,
		ServerInfo: ServerInfo{},
	}
	// viper自动读取当前环境变量？
	conf.viper.AutomaticEnv()
	conf.viper.AllowEmptyEnv(true)
	// viper 加载指定配置文件
	conf.viper.SetConfigFile(configFile)
	//conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	if configFile == "" {
		workDir, _ := os.Getwd()
		conf.viper.AddConfigPath(workDir + "/config")
		conf.viper.AddConfigPath("/etc/project/")
	}

	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln("配置初始化失败", err)
	}
	// viper打开热更新配置功能，监听配置文件修改（应当放在main函数整体监听后重新初始化）
	//conf.viper.WatchConfig()
	//conf.viper.OnConfigChange(func(e fsnotify.Event) {
	//	log.Println("Config file changed:", e.Name)
	//})
	if err := conf.viper.Unmarshal(conf); err != nil {
		log.Fatalln("服务信息初始化失败：", err)
	}
	return conf
}

func (c *Config) InitZapLog() {
	//从配置中读取日志配置，初始化日志
	// 这个是手动配置日志配置对象;如果config.yaml文件有层级的话不可以用一层struct平铺展开，只能按照config层级嵌套。太麻烦
	logConfig := &common.LogConfig{
		Level:            c.Mode,
		InfoLogFilename:  c.viper.GetString("log.infoLogFilename"),
		WarnLogFilename:  c.viper.GetString("log.warnLogFilename"),
		ErrorLogFilename: c.viper.GetString("log.errorLogFilename"),
		MaxSize:          c.viper.GetInt("log.maxSize"),
		MaxAge:           c.viper.GetInt("log.maxAge"),
		MaxBackups:       c.viper.GetInt("log.maxBackupFiles"),
	}
	err := common.InitLogger(logConfig)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) InitRedisDB() {
	redisInfo := &model.RedisInfo{
		Host:     c.viper.GetString("redis.host"),
		Port:     c.viper.GetInt("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
	model.InitRedisCache(redisInfo)
}
