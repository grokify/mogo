package slogutil

import (
	"context"
	"io"
	"log/slog"
)

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
