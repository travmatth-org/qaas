package logger

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/travmatth-org/qaas/internal/types"
)

// var instance *zerolog.Logger
var instance *types.ZLog

var destination io.Writer = os.Stdout

// SetLogger sets the logger instance
func SetLogger(base *types.ZLog) {
	instance = base
}

// SetDestination sets the output location of the logger
func SetDestination(dest io.Writer) {
	destination = dest
}

// GetLogger gets the logger instance, or configures one if nil
func GetLogger() *types.ZLog {
	if instance != nil {
		return instance
	}
	zerolog.TimeFieldFormat = time.RFC1123
	base := zerolog.New(destination)
	base = base.With().Timestamp().Caller().Str("role", "qaas").Logger()
	SetLogger(&base)
	return instance
}

// Error returns a zerolog Error logger
func Error() *types.ZLEvent {
	return GetLogger().Error()
}

// Info returns a zerolog Info logger
func Info() *types.ZLEvent {
	return GetLogger().Info()
}

// Warn returns a zerolog Warn logger
func Warn() *types.ZLEvent {
	return GetLogger().Warn()
}

// Debug returns a zerolog Debug logger
func Debug() *types.ZLEvent {
	return GetLogger().Debug()
}

// ErrorReq returns a zerolog Error logger from *http.Request
func ErrorReq(r *http.Request) *types.ZLEvent {
	return hlog.FromRequest(r).Error()
}

// InfoReq returns a zerolog Info logger from *http.Request
func InfoReq(r *http.Request) *types.ZLEvent {
	return hlog.FromRequest(r).Info()
}

// WarnReq returns a zerolog Warn logger from *http.Request
func WarnReq(r *http.Request) *types.ZLEvent {
	return hlog.FromRequest(r).Warn()
}

// DebugReq returns a zerolog Debug logger from *http.Request
func DebugReq(r *http.Request) *types.ZLEvent {
	return hlog.FromRequest(r).Debug()
}

// InfoOrErr returns an info or error logger depending on err status
func InfoOrErr(err error) *types.ZLEvent {
	if err != nil {
		return Error().Err(err)
	}
	return Info()
}
