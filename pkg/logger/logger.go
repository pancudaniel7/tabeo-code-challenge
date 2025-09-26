package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

var LevelTrace = slog.Level(-8)

type Config struct {
	Level string
	JSON  bool
}

func NewDefaultLogger(cfg Config) AppLogger {
	var h slog.Handler
	lvl := parseLevel(cfg.Level)
	if cfg.JSON {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl, AddSource: true})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl, AddSource: true})
	}
	return &slogAdapter{base: slog.New(h)}
}

type slogAdapter struct {
	base *slog.Logger
}

func (a *slogAdapter) Trace(ctx context.Context, msg string, kv ...any) {
	a.base.Log(ctx, LevelTrace, msg, kv...)
}
func (a *slogAdapter) Debug(ctx context.Context, msg string, kv ...any) {
	a.base.DebugContext(ctx, msg, kv...)
}
func (a *slogAdapter) Info(ctx context.Context, msg string, kv ...any) {
	a.base.InfoContext(ctx, msg, kv...)
}
func (a *slogAdapter) Warn(ctx context.Context, msg string, kv ...any) {
	a.base.WarnContext(ctx, msg, kv...)
}
func (a *slogAdapter) Error(ctx context.Context, err error, msg string, kv ...any) {
	a.base.ErrorContext(ctx, msg, append(kv, "error", err)...)
}
func (a *slogAdapter) Fatal(ctx context.Context, err error, msg string, kv ...any) {
	a.base.ErrorContext(ctx, msg, append(kv, "error", err)...)
	os.Exit(1)
}
func (a *slogAdapter) With(kv ...any) AppLogger {
	return &slogAdapter{base: a.base.With(kv...)}
}

func parseLevel(s string) slog.Leveler {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "trace":
		return LevelTrace
	case "debug":
		return slog.LevelDebug
	case "info", "":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
