package api

import (
	"net/http"

	"github.com/utsabbera/task-master/pkg/middleware"
)

func NewRouter(handler Handler, middlewares ...middleware.Middleware) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /tasks", handler.Create)
	router.HandleFunc("GET /tasks", handler.List)
	router.HandleFunc("GET /tasks/{id}", handler.Get)
	router.HandleFunc("PUT /tasks/{id}", handler.Update)
	router.HandleFunc("DELETE /tasks/{id}", handler.Delete)

	return middleware.Bind(router, middlewares...)
}
