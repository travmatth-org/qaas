package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	n := len(middlewares) - 1
	return func(h http.Handler) http.Handler {
		for i := range middlewares {
			h = middlewares[n-i](h)
		}
		return h
	}
}
