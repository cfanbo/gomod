package log

import (
	"context"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func Init() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func SetDebugLevel(flag bool) {
	if flag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}

func With() zerolog.Context {
	return zlog.With()
}

func Err(err error) *zerolog.Event {
	return zlog.Err(err)
}

func Trace() *zerolog.Event {
	return zlog.Trace()
}

func Debug() *zerolog.Event {
	return zlog.Debug()
}

func Info() *zerolog.Event {
	return zlog.Info()
}

func Warn() *zerolog.Event {
	return zlog.Warn()
}

func Error() *zerolog.Event {
	return zlog.Error()
}

func WithLevel(level zerolog.Level) *zerolog.Event {
	return zlog.WithLevel(level)
}

func Print(v ...interface{}) {
	zlog.Print(v...)
}

func Printf(format string, v ...interface{}) {
	zlog.Printf(format, v...)
}

func Ctx(ctx context.Context) *zerolog.Logger {
	return zlog.Ctx(ctx)
}
