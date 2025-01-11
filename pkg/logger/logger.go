package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log переменная для доступа к логгеру
var Log = zap.NewNop()

// Initialize метод инициализации логгера
func Initialize(level string, output []string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if output != nil {
		cfg.OutputPaths = output
	}
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}
