// internal/logger/logger.go
package logger

import (
	"go.uber.org/zap"
)

func New(level string) (*zap.Logger, error) {
	var cfg zap.Config
	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	return cfg.Build()
}
