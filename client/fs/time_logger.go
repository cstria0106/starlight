package fs

import (
	"context"
	"fmt"
	"io"
	"os"
)

type accessLoggerKey struct{}

type AccessLogger struct {
	file io.Writer
}

func WithAccessLogger(ctx context.Context, filepath string) context.Context {
	var file io.Writer
	var err error
	if file, err = os.Create(filepath); err != nil {
		panic(fmt.Sprintf("file to create access logger file on %s", filepath))
	}

	return context.WithValue(ctx, accessLoggerKey{}, &AccessLogger{
		file: file,
	})
}

func GetAccessLogger(ctx context.Context) *AccessLogger {
	l := ctx.Value(accessLoggerKey{})
	if l == nil {
		panic("No access logger in context")
	}

	return l.(*AccessLogger)
}
