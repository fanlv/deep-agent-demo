package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

var defaultLogger *slog.Logger

func init() {
	defaultLogger = slog.New(&plainHandler{w: os.Stdout, level: slog.LevelInfo})
}

// plainHandler outputs logs in a simple plain-text format that preserves newlines in messages.
type plainHandler struct {
	w     io.Writer
	level slog.Level
}

func (h *plainHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *plainHandler) Handle(_ context.Context, r slog.Record) error {
	ts := r.Time.Format(time.DateTime)
	_, err := fmt.Fprintf(h.w, "%s %s %s\n", ts, r.Level.String(), r.Message)
	return err
}

func (h *plainHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *plainHandler) WithGroup(_ string) slog.Handler      { return h }

func SetDefault(l *slog.Logger) {
	defaultLogger = l
}

func Default() *slog.Logger {
	return defaultLogger
}

func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
	os.Exit(1)
}

func Debugf(ctx context.Context, format string, args ...any) {
	defaultLogger.DebugContext(ctx, fmt.Sprintf(format, args...))
}

func Infof(ctx context.Context, format string, args ...any) {
	defaultLogger.InfoContext(ctx, fmt.Sprintf(format, args...))
}

func Warnf(ctx context.Context, format string, args ...any) {
	defaultLogger.WarnContext(ctx, fmt.Sprintf(format, args...))
}

func Errorf(ctx context.Context, format string, args ...any) {
	defaultLogger.ErrorContext(ctx, fmt.Sprintf(format, args...))
}

func Fatalf(ctx context.Context, format string, args ...any) {
	defaultLogger.ErrorContext(ctx, fmt.Sprintf(format, args...))
	os.Exit(1)
}
