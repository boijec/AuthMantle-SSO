package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

const LoggerContextKey ContextKey = "logger"
const ResponseKey = "response"
const DurationKey = "duration"

type ResponseWrapper struct {
	http.ResponseWriter
	status int
}

func (rw *ResponseWrapper) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	rw.status = status
}

func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		corrId := r.Header.Get("X-Correlation-ID")
		if corrId == "" {
			corrId = uuid.NewString()
		}
		logger := slog.With(
			"correlation_id", corrId,
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"path", r.URL.Path,
		)
		ctx = context.WithValue(ctx, LoggerContextKey, logger)
		r = r.WithContext(ctx)

		logger.DebugContext(ctx, "Incoming request")
		s := time.Now()
		responseWrapper := &ResponseWrapper{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		defer func() {
			logger.InfoContext(ctx, "Finished request", ResponseKey, responseWrapper.status, DurationKey, time.Since(s))
		}()
		next.ServeHTTP(responseWrapper, r)
		// TODO change to correlation id MW, when the request throws an error, this is dumb..

	})
}
