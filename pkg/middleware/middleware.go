package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// REVIEW: This can be a generic method to bind handlers.

// Bind applies a series of middlewares to an HTTP handler.
func Bind(handler http.Handler, middlewares ...Middleware) http.Handler {
	router := handler

	if len(middlewares) == 0 {
		return handler
	}

	for _, mw := range middlewares {
		router = mw(router)
	}

	return router
}
