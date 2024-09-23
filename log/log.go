package log

import (
	"flow/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// InitLogger 初始化全局Zap日志
func InitLogger() {
	// 设置日志分割配置
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.XConfig.Log.Path,
		MaxSize:    config.XConfig.Log.MaxSize, // megabytes
		MaxBackups: config.XConfig.Log.MaxBackups,
		MaxAge:     config.XConfig.Log.MaxAge,   // days
		Compress:   config.XConfig.Log.Compress, // 是否压缩
	}

	// 创建 Zap 核心
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(lumberjackLogger),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			level, err := zap.ParseAtomicLevel(config.XConfig.Log.LevelFile)
			if err != nil {
				panic(err)
			}
			return lvl >= level.Level()
		}),
	)

	// 创建控制台核心
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			level, err := zap.ParseAtomicLevel(config.XConfig.Log.LevelConsole)
			if err != nil {
				panic(err)
			}
			return lvl >= level.Level()
		}),
	)

	// 合并核心
	core := zapcore.NewTee(fileCore, consoleCore)

	// 创建 Logger
	logger := zap.New(core)

	// 将当前日志记录器设置为全局默认记录器
	zap.ReplaceGlobals(logger)

}
