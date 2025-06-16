package api

import (
	"net/http"

	swagger "github.com/swaggo/http-swagger"
	"github.com/utsabbera/task-master/pkg/middleware"
)

func newTaskRouter(handler Handler, middlewares ...middleware.Middleware) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /", handler.Create)
	router.HandleFunc("GET /", handler.List)
	router.HandleFunc("GET /{id}", handler.Get)
	router.HandleFunc("PUT /{id}", handler.Update)
	router.HandleFunc("DELETE /{id}", handler.Delete)

	return middleware.Bind(router, middlewares...)
}

func NewRouter(handler Handler, middlewares ...middleware.Middleware) http.Handler {

	router := http.NewServeMux()
	router.Handle("/tasks/", newTaskRouter(handler, middlewares...))
	router.Handle("/swagger/", swagger.WrapHandler)

	return router
}
