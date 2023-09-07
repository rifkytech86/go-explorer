package logger

import (
	"context"
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
	ctx    context.Context
}

func NewZapLogger(loggerType string, ctx context.Context) *ZapLogger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableStacktrace = true
	zapLogger, _ := zapConfig.Build(zap.AddCallerSkip(2))
	defer zapLogger.Sync()
	return &ZapLogger{logger: zapLogger, ctx: ctx}
}

func (l *ZapLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.Debug(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Info(msg string, fields map[string]interface{}) {

	l.logger.Info(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Warn(msg string, fields map[string]interface{}) {

	l.logger.Warn(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Error(msg string, fields map[string]interface{}) {

	l.logger.Error(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Fatal(msg string, fields map[string]interface{}) {

	l.logger.Fatal(msg, zap.Any("args", fields))
}

func (l *ZapLogger) addContextCommonFields(fields map[string]interface{}) {
	if l.ctx != nil {
		for k, v := range l.ctx.Value("commonFields").(map[string]interface{}) {
			if _, ok := fields[k]; !ok {
				fields[k] = v
			}
		}
	}
}
