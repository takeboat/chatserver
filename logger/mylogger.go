package logger

import "log/slog"

type Logger struct {
	logger *slog.Logger
}

func NewLogger(name string) *Logger {
	return &Logger{
		logger: slog.New(slog.NewJSONHandler(nil, nil)),
	}
}
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		logger: l.logger.WithGroup(name),
	}
}