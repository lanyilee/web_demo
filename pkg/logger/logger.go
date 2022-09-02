package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var levels = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLevel(level string) zapcore.Level {
	if level, ok := levels[level]; ok {
		return level
	} else {
		return zapcore.DebugLevel
	}
}

func Setup(level string, logJSON bool) {
	// 创建控制台Core记录器接口
	consoleCore := zapcore.NewCore(
		getConsoleEncoder(logJSON),
		zapcore.Lock(os.Stdout),
		getLevel(level),
	)
	cores := []zapcore.Core{consoleCore}
	// Option可选配置
	options := []zap.Option{
		zap.AddCaller(),                   // Caller调用显示
		zap.AddStacktrace(zap.ErrorLevel), // 堆栈跟踪级别
		//zap.Fields(zap.String("label", "WeBase-Zap")), // 全局添加字段
	}

	// 构建一个 zap 实例
	logger := zap.New(
		zapcore.NewTee(cores...),
		options...,
	)

	zap.ReplaceGlobals(logger) // 设为全局zap实例

}

// getConsoleEncoder 控制台日志格式
func getConsoleEncoder(logJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	if logJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}
