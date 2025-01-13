package log

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

type Logger struct {
	*slog.Logger
}

func Default() *Logger {
	return &Logger{slog.Default()}
}

func (l *Logger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if !l.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{l.Logger.With(args...)}
}

func (l *Logger) DebugCtx(ctx context.Context, format string, args ...any) {
	l.log(ctx, slog.LevelDebug, fmt.Sprintf(format, args...))
}
func (l *Logger) InfoCtx(ctx context.Context, format string, args ...any) {
	l.log(ctx, slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (l *Logger) WarnCtx(ctx context.Context, format string, args ...any) {
	l.log(ctx, slog.LevelWarn, fmt.Sprintf(format, args...))
}

func (l *Logger) ErrorCtx(ctx context.Context, format string, args ...any) {
	l.log(ctx, slog.LevelError, fmt.Sprintf(format, args...))
}

func (l *Logger) Debug(format string, args ...any) {
	l.log(context.Background(), slog.LevelDebug, fmt.Sprintf(format, args...))
}
func (l *Logger) Info(format string, args ...any) {
	l.log(context.Background(), slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(format string, args ...any) {
	l.log(context.Background(), slog.LevelWarn, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(format string, args ...any) {
	l.log(context.Background(), slog.LevelError, fmt.Sprintf(format, args...))
}

func With(args ...any) *Logger {
	return Default().With(args...)
}

func Debug(format string, args ...any) {
	Default().log(context.Background(), slog.LevelDebug, fmt.Sprintf(format, args...))
}
func Info(format string, args ...any) {
	Default().log(context.Background(), slog.LevelInfo, fmt.Sprintf(format, args...))
}

func Warn(format string, args ...any) {
	Default().log(context.Background(), slog.LevelWarn, fmt.Sprintf(format, args...))
}

func Error(format string, args ...any) {
	Default().log(context.Background(), slog.LevelError, fmt.Sprintf(format, args...))
}

func DebugCtx(ctx context.Context, format string, args ...any) {
	Default().log(ctx, slog.LevelDebug, fmt.Sprintf(format, args...))
}
func InfoCtx(ctx context.Context, format string, args ...any) {
	Default().log(ctx, slog.LevelInfo, fmt.Sprintf(format, args...))
}

func WarnCtx(ctx context.Context, format string, args ...any) {
	Default().log(ctx, slog.LevelWarn, fmt.Sprintf(format, args...))
}

func ErrorCtx(ctx context.Context, format string, args ...any) {
	Default().log(ctx, slog.LevelError, fmt.Sprintf(format, args...))
}
