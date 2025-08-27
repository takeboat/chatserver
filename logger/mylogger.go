package logger

import (
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
	level  slog.LevelVar
}

type LoggerOption func(*Logger)

// WithGroup 添加组
func WithGroup(name string) LoggerOption {
	return func(l *Logger) {
		l.logger = l.logger.WithGroup(name)
	}
}

// WithLevel 设置日志级别
func WithLevel(level slog.Level) LoggerOption {
	return func(l *Logger) {
		l.level.Set(level)
	}
}


// WithName 设置组件名称（component）
func WithName(name string) LoggerOption {
	return func(l *Logger) {
		l.logger = l.logger.With("component", name)
	}
}

// NewLogger 创建一个新的 Logger 实例
func NewLogger(opts ...LoggerOption) *Logger {
	l := &Logger{
		logger: slog.Default(),
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Debug 记录 Debug 级别的日志
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info 记录 Info 级别的日志
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn 记录 Warn 级别的日志
func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error 记录 Error 级别的日志
func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

