package middleware

import (
	"log/slog"
	"net/http"
)

type MiddlewareFunc func(next http.Handler) http.Handler

func NewSlogMiddleware(logger *slog.Logger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("request", slog.String("method", r.Method), slog.String("uri", r.RequestURI), slog.String("remoteAddr", r.RemoteAddr))
			next.ServeHTTP(w, r)
		})
	}
}
