package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/app/provider"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logg := provider.CalendarContainer.GetLogger()
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		latency := time.Since(start)

		ctx := context.WithValue(r.Context(), logg.GetContextKey("client_ip"), strings.Split(r.RemoteAddr, ":")[0])
		ctx = context.WithValue(ctx, logg.GetContextKey("http_method"), r.Method)
		ctx = context.WithValue(ctx, logg.GetContextKey("http_path"), r.URL.Path)
		ctx = context.WithValue(ctx, logg.GetContextKey("http_version"), r.Proto)
		ctx = context.WithValue(ctx, logg.GetContextKey("http_response_code"), fmt.Sprintf("%d", rw.statusCode))
		ctx = context.WithValue(ctx, logg.GetContextKey("latency"), latency.String())
		ctx = context.WithValue(ctx, logg.GetContextKey("user_agent"), r.UserAgent())

		logg.Info(ctx, "request log")
	})
}
