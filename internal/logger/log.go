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

var destination io.Writer = os.Stdout

// SetLogger sets the logger instance
func SetLogger(base *zerolog.Logger) {
	instance = base
}

// SetLogger sets the output location of the logger
func SetDestination(dest io.Writer) {
	destination = dest
}

// GetLogger gets the logger instance, or configures one if nil
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

// Error returns a zerolog Error logger
func Error() *zerolog.Event {
	return GetLogger().Error()
}

// Info returns a zerolog Info logger
func Info() *zerolog.Event {
	return GetLogger().Info()
}

// Warn returns a zerolog Warn logger
func Warn() *zerolog.Event {
	return GetLogger().Warn()
}

// Debug returns a zerolog Debug logger
func Debug() *zerolog.Event {
	return GetLogger().Debug()
}

// ErrorReq returns a zerolog Error logger from *http.Request
func ErrorReq(r *http.Request) *zerolog.Event {
	return hlog.FromRequest(r).Error()
}

// InfoReq returns a zerolog Info logger from *http.Request
func InfoReq(r *http.Request) *zerolog.Event {
	return hlog.FromRequest(r).Info()
}

// WarnReq returns a zerolog Warn logger from *http.Request
func WarnReq(r *http.Request) *zerolog.Event {
	return hlog.FromRequest(r).Warn()
}

// DebugReq returns a zerolog Debug logger from *http.Request
func DebugReq(r *http.Request) *zerolog.Event {
	return hlog.FromRequest(r).Debug()
}
