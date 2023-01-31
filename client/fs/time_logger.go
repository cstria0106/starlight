package fs

import "context"

type timeLoggerKey struct{}

type TimeLogger struct{}

func WithTimeLogger(ctx context.Context, filepath string) context.Context {}

func GetTimeLogger(ctx context.Context) *TimeLogger {
	l := ctx.Value(timeLoggerKey{})
	if l == nil {
		panic("No time logger in context")
	}

	return l.(*TimeLogger)
}
