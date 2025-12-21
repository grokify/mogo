package slogutil

import (
	"context"
	"io"
	"log/slog"
)

// loggerKey is the context key for storing a request-scoped logger
type loggerKey struct{}

// ContextWithLogger returns a new context with the given logger attached.
// Use this to pass request-scoped loggers (with trace IDs, user IDs, etc.)
// that will be used for logging within that request.
func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// LoggerFromContext returns the logger from context if present,
// otherwise returns the fallback logger.
func LoggerFromContext(ctx context.Context, fallback *slog.Logger) *slog.Logger {
	if ctx == nil {
		return fallback
	}
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok && logger != nil {
		return logger
	}
	return fallback
}

type NullHandler struct{}

func (NullHandler) Enabled(_ context.Context, _ slog.Level) bool  { return false }
func (NullHandler) Handle(_ context.Context, _ slog.Record) error { return nil }
func (h NullHandler) WithAttrs(_ []slog.Attr) slog.Handler        { return h }
func (h NullHandler) WithGroup(_ string) slog.Handler             { return h }

// Null is an inexpensive nop logger.
func Null() *slog.Logger {
	return slog.New(NullHandler{})
}

// Discard is an expensive nop logger using the stdlib.
func Discard() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}
