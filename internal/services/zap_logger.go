package services

import (
	"github.com/coci/cutme/internal/core/ports"
	"go.uber.org/zap"
)

type ZapLogger struct {
	ports.Logger

	logger *zap.Logger
}

func NewZapLogger() *ZapLogger {
	logger, _ := zap.NewProduction()
	return &ZapLogger{
		logger: logger,
	}
}

func (z *ZapLogger) Info(msg string, fields ...ports.Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	z.logger.Info(msg, zapFields...)
}

func (z *ZapLogger) Warn(msg string, fields ...ports.Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	z.logger.Warn(msg, zapFields...)
}

func (z *ZapLogger) Error(msg string, fields ...ports.Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	z.logger.Error(msg, zapFields...)
}

func (z *ZapLogger) Fatal(msg string, fields ...ports.Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	z.logger.Fatal(msg, zapFields...)
}
