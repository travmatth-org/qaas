package handlers

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
	"github.com/travmatth-org/qaas/internal/afs"
	"github.com/travmatth-org/qaas/internal/api"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/logger"
)

// Handler manages the endpoints of the Server
type Handler struct {
	api    *api.API
	fs     *afs.AFS
	tables config.Tables
}

// Opt are the configuratio nfunctions used set options on the hander struct
type Opt func(h *Handler) (*Handler, error)

// RouteMap returns a map linking endpoint endpoint names to their handler
func (h *Handler) RouteMap() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/":       Static(h.fs.Use("index")),
		"/get":    h.Get,
		"/put":    h.Put,
		"/random": h.Random,
		"/404":    Static(h.fs.Use("404"))}
}

// New constructs and returns a Handler struct with the specified opts
func New(opts ...Opt) (*Handler, error) {
	var (
		h         = &Handler{}
		err error = nil
	)
	for _, opt := range opts {
		h, err = opt(h)
		if err != nil {
			return nil, err
		}
	}
	return h, err
}

// WithAPI sets an API for the handler struct
func WithAPI(a *api.API) Opt {
	return func(h *Handler) (*Handler, error) {
		h.api = a
		return h, nil
	}
}

// WithFS inserts the given file system client into the server
func WithFS(fs *afs.AFS, static string) Opt {
	return func(h *Handler) (*Handler, error) {
		h.fs = fs
		return h, fs.LoadAssets(static)
	}
}

// Route composes endpoints by wrapping destination handler with handler
// pipeline providing tracing with aws x-ray, injecting logging middleware
// with request details into the context, and error recovery middleware,
// and gzipping the response
func Route(h http.HandlerFunc, isProd bool) http.HandlerFunc {
	handler := alice.New(
		Recover,
		hlog.NewHandler(*logger.GetLogger()),
		hlog.RequestIDHandler("req_id", "Request-Id"),
		hlog.RemoteAddrHandler("ip"),
		hlog.RequestHandler("dest"),
		hlog.RefererHandler("referer"),
		gziphandler.GzipHandler,
		Log,
	).ThenFunc(h)
	if isProd {
		namer := xray.NewFixedSegmentNamer("qaas-httpd")
		return xray.Handler(namer, handler).ServeHTTP
	}
	return handler.ServeHTTP
}
