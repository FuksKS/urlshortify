package logger

import (
	"fmt"
	"go.uber.org/zap"
)

const (
	LoggerLevelINFO  = "INFO"
	LoggerLevelDEBUG = "DEBUG"
)

var Log *zap.Logger = zap.NewNop()

// Init инициализирует синглтон логера с необходимым уровнем логирования.
func Init(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("logger-Init-ParseAtomicLevel-err: %w", err)
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("logger-Init-cfg.Build-err: %w", err)
	}

	Log = zl
	return nil
}
