package slogutil

import (
	"context"
	"fmt"
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

func LogOrNotAny(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, args ...any) {
	attrs := AttrsFromAny(args...)
	LogOrNot(ctx, logger, level, msg, attrs...)
}

func AttrsFromAny(args ...any) []slog.Attr {
	var attrs []slog.Attr
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", key)
		}
		if i < len(args)-1 {
			attrs = append(attrs, slog.Any(key, args[i+1]))
		} else {
			attrs = append(attrs, slog.Any(key, nil))
		}
	}
	return attrs
}
