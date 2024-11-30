package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Encoding = "console"     // Verbose output, JSON is not good for your eyes...
	config.DisableStacktrace = true // Turn off stacktrace print

	baseLogger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1)) // Skip one call from stacktrace
	if err != nil {
		panic(err)
	}

	logger = baseLogger
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func ErrorWithStacktrace(msg string, fields ...zap.Field) {
	logger.Error(msg, append(fields, zap.Stack("stacktrace"))...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, append(fields, zap.Stack("stacktrace"))...)
}

func Sync() {
	_ = logger.Sync() // Close logger
}
