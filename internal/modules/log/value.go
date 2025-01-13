package log

import "log/slog"

func ErrorString(err error) slog.Attr {
	return slog.String("err", err.Error())
}

var (
	Any      = slog.Any
	String   = slog.String
	Int      = slog.Int
	Int64    = slog.Int64
	Uint64   = slog.Uint64
	Float64  = slog.Float64
	Duration = slog.Duration
	Bool     = slog.Bool
	Time     = slog.Time
)
