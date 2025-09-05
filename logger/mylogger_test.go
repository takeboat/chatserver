package logger

import (
	"log/slog"
	"testing"
)

func TestMyLogger(t *testing.T) {
	logger := NewLogger(WithGroup("test_root"))
	logger.SetLevel(slog.LevelDebug)
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	newlogger := logger.WithGroup("test_group")
	newlogger.Debug("debug")
	newlogger.Debug("debug")
	newlogger.Info("info")
	newlogger.Warn("warn")
	newlogger.Error("error")
}
