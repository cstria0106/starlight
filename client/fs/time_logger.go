package fs

import (
	"context"
	"fmt"
	"os"
	"time"
)

type accessLoggerKey struct{}

type AccessLogger struct {
	file *os.File
}

func (l *AccessLogger) Log(file string) {
	t := time.Now()
	l.file.WriteString(fmt.Sprintf("%s\t%s\n", t.Format("12:00:00.00"), file))
}

func WithAccessLogger(ctx context.Context, filepath string) context.Context {
	var file *os.File
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
