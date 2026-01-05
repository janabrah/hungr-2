package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
)

var Log *slog.Logger

func Init() {
	env := os.Getenv("ENV")
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	Log = slog.New(handler)
}

type contextKey string

const requestIDKey contextKey = "request_id"

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

func Error(ctx context.Context, msg string, err error, args ...any) {
	stack := make([]byte, 4096)
	n := runtime.Stack(stack, false)
	allArgs := append([]any{
		"error", err.Error(),
		"request_id", GetRequestID(ctx),
		"stack", string(stack[:n]),
	}, args...)
	Log.Error(msg, allArgs...)
}

func Info(ctx context.Context, msg string, args ...any) {
	allArgs := append([]any{"request_id", GetRequestID(ctx)}, args...)
	Log.Info(msg, allArgs...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	allArgs := append([]any{"request_id", GetRequestID(ctx)}, args...)
	Log.Debug(msg, allArgs...)
}
