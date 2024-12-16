package logging

import (
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	w    http.ResponseWriter
	Code int
}

func newLoggingResponseWriter(ww http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		w:    ww,
		Code: 0,
	}
}

func (w *loggingResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.Code = statusCode
	w.w.WriteHeader(statusCode)
}

func NewLoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w2 := newLoggingResponseWriter(w)
		start := time.Now()
		handler.ServeHTTP(w2, r)
		slog.Info("incoming request", "mehod", r.Method, "url", r.URL.Path, "status", w2.Code, "duration", time.Since(start).String())
	})
}
