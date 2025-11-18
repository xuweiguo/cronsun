package bootstrap

import (
	"fmt"
	slog "log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"cronsun/log"
)

func initLogger(level int) error {
	lcf := zap.NewDevelopmentConfig()
	lcf.Level.SetLevel(zapcore.Level(level))
	lcf.Development = false
	lcf.DisableStacktrace = true
	logger, err := lcf.Build(zap.AddCallerSkip(1))
	if err != nil {
		return fmt.Errorf("new log err: %w", err)
	}
	log.SetLogger(logger.Sugar())
	slog.SetFlags(slog.LstdFlags | slog.Lshortfile)
	return nil
}
