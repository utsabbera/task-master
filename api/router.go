package api

import (
	"net/http"

	swagger "github.com/swaggo/http-swagger"
	"github.com/utsabbera/task-master/pkg/middleware"
)

// NewTaskRouter creates a new HTTP router for task-related endpoints.
func NewTaskRouter(handler Handler, middlewares ...middleware.Middleware) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /tasks", handler.Create)
	router.HandleFunc("GET /tasks", handler.List)
	router.HandleFunc("GET /tasks/{id}", handler.Get)
	router.HandleFunc("PUT /tasks/{id}", handler.Update)
	router.HandleFunc("DELETE /tasks/{id}", handler.Delete)

	return middleware.Bind(router, middlewares...)
}

// NewRouter creates the main HTTP router for the API.
func NewRouter(handler Handler, middlewares ...middleware.Middleware) http.Handler {

	router := http.NewServeMux()
	router.Handle("/", NewTaskRouter(handler, middlewares...))
	router.Handle("/swagger/", swagger.WrapHandler)

	return router
}
