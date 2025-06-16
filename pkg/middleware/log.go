package middleware

import (
	"log"
	"net/http"
	"time"
)

func Log() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)
			duration := time.Since(start)

			log.Printf(
				"%s %s %s %d %s %v",
				r.RemoteAddr,
				r.Method,
				r.URL.Path,
				rw.statusCode,
				http.StatusText(rw.statusCode),
				duration,
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
