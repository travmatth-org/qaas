package logger

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var instance *zerolog.Logger
var destination io.Writer = os.Stderr

func SetLogger(base *zerolog.Logger) {
	instance = base
}

func GetLogger() *zerolog.Logger {
	if instance != nil {
		return instance
	}

	zerolog.TimeFieldFormat = time.RFC1123
	base := zerolog.New(destination).With().
		Timestamp().
		Caller().
		Str("role", "faas").
		Logger()
	SetLogger(&base)
	return instance
}

// Error
func Error() *zerolog.Event {
	return GetLogger().Error()
}

func ErrorReq(r *http.Request) *zerolog.Event {
	return hlog.FromRequest(r).Error()
}

// Info

func Info() *zerolog.Event {
	return GetLogger().Info()
}

func InfoReq(r *http.Request) *zerolog.Event {
	return hlog.FromRequest(r).Info()
}

// Warn

func Warn() *zerolog.Event {
	return GetLogger().Warn()
}

func WarnReq(r *http.Request) *zerolog.Event {
	return hlog.FromRequest(r).Warn()
}

// Debug

func Debug() *zerolog.Event {
	return GetLogger().Debug()
}

func DebugReq(r *http.Request) *zerolog.Event {
	return hlog.FromRequest(r).Debug()
}
