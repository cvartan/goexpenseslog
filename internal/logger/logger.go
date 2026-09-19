package logger

import (
	"context"
	"log/slog"
)

func GetLogger(ctx context.Context) *slog.Logger {
	l := ctx.Value("logger")
	if l == nil {
		panic("logger is not defined in this context")
	}
	logger, ok := l.(*slog.Logger)
	if !ok {
		panic("returned value from context for logger is not *slog.Logger")
	}
	return logger
}
