package logger

import (
	"context"
	"go.uber.org/zap"
)

var globalLogger *zap.Logger

func InitLogger() error {
	var err error
	globalLogger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

type Log struct {
	logger *zap.Logger
	ctx    context.Context
}

func GetLogger() *zap.Logger {
	return globalLogger
}

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

type LoggerWrapper struct {
	logger Logger
}

var sharedLogger Logger

// SetSharedLogger sets the shared logger instance
func SetSharedLogger(logger Logger) {
	sharedLogger = logger
}

func NewLoggerWrapper(loggerType string, ctx context.Context) Logger {
	logger := NewZapLogger(loggerType, ctx)
	SetSharedLogger(logger)
	return logger
}

func Debug(msg string, fields map[string]interface{}) {
	if sharedLogger != nil {
		sharedLogger.Debug(msg, fields)
	}

}

func Info(msg string, fields map[string]interface{}) {
	if sharedLogger != nil {
		sharedLogger.Info(msg, fields)
	}
}

func Warn(msg string, fields map[string]interface{}) {
	if sharedLogger != nil {
		sharedLogger.Warn(msg, fields)
	}

}

func Error(msg string, fields map[string]interface{}) {
	if sharedLogger != nil {
		sharedLogger.Error(msg, fields)
	}
}

func Fatal(msg string, fields map[string]interface{}) {
	if sharedLogger != nil {
		sharedLogger.Fatal(msg, fields)
	}
}
