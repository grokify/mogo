package slogutil

import (
	"context"
	"log/slog"
)

func LogOrNot(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, attrs ...slog.Attr) {
	if logger == nil {
		return
	} else if ctx == nil {
		ctx = context.Background()
	}
	logger.LogAttrs(ctx, level, msg, attrs...)
}
