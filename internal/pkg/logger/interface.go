package logger

import "context"

type AppLogger interface {
	Trace(ctx context.Context, msg string, kv ...any)
	Debug(ctx context.Context, msg string, kv ...any)
	Info(ctx context.Context, msg string, kv ...any)
	Warn(ctx context.Context, msg string, kv ...any)
	Error(ctx context.Context, err error, msg string, kv ...any)
	Fatal(ctx context.Context, err error, msg string, kv ...any)
	With(kv ...any) AppLogger
}
