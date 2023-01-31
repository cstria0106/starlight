package fs

import (
	"context"
	"fmt"
	"io"
	"os"
)

type timeLoggerKey struct{}

type TimeLogger struct {
	file io.Writer
}

func WithTimeLogger(ctx context.Context, filepath string) context.Context {
	var file io.Writer
	var err error
	if file, err = os.Create(filepath); err != nil {
		panic(fmt.Sprintf("file to create time logger file on %s", filepath))
	}

	return context.WithValue(ctx, timeLoggerKey{}, &TimeLogger{
		file: file,
	})
}

func GetTimeLogger(ctx context.Context) *TimeLogger {
	l := ctx.Value(timeLoggerKey{})
	if l == nil {
		panic("No time logger in context")
	}

	return l.(*TimeLogger)
}
