// internal/logger/logger.go
package logger

import (
	"go.uber.org/zap"
)

func NewLogger(level string) *zap.Logger {
	var cfg zap.Config
	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	logger, _ := cfg.Build()
	return logger
}
