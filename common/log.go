package common

import (
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// Logger 此文件在于实现gin框架中引入zap日志库，实现日志的级别划分、标准化日志格式以及日志轮转
var Logger *zap.Logger

type LogConfig struct {
	Level            string `json:"level"`
	InfoLogFilename  string `json:"info_log_filename"`
	WarnLogFilename  string `json:"warn_log_filename"`
	ErrorLogFilename string `json:"error_log_filename"`
	MaxSize          int    `json:"maxsize"`
	MaxAge           int    `json:"max_age"`
	MaxBackups       int    `json:"max_backups"`
}

// InitLogger 通过日志相关配置初始化Logger
func InitLogger(cfg *LogConfig) (err error) {
	writeSyncerInfo := getLogWriter(cfg.InfoLogFilename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	writeSyncerWarn := getLogWriter(cfg.WarnLogFilename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	writeSyncerError := getLogWriter(cfg.ErrorLogFilename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	encoder := getEncoder()
	//文件输出
	infoCore := zapcore.NewCore(encoder, writeSyncerInfo, zapcore.DebugLevel)
	warnCore := zapcore.NewCore(encoder, writeSyncerWarn, zapcore.WarnLevel)
	debugCore := zapcore.NewCore(encoder, writeSyncerError, zapcore.ErrorLevel)
	//标准输出
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	var std zapcore.Core
	switch cfg.Level {
	case "dev":
		std = zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
	case "uat":
		std = zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.WarnLevel)
	case "prod":
		std = zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.ErrorLevel)
	default:
		std = zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
	}
	core := zapcore.NewTee(debugCore, infoCore, warnCore, std)
	Logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	// 日志分割
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		Logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					Logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
