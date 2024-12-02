package logger

import (
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Module = fx.Provide(
	NewLogger,
)

func NewLogger(cfg *viper.Viper) *zap.Logger {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.DisableStacktrace = true

	if cfg.GetBool("app.debug") {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.DisableStacktrace = false
		config.DisableCaller = false
	}

	logger, _ := config.Build()
	return logger
}

// NewLID create loggerID
func NewLID() zap.Field {
	return zap.String("lid", uuid.NewString())
}
