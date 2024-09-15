package middleware_test

import (
	"authmantle-sso/middleware"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"net/http"
	"testing"
)

func BenchmarkRegisterMiddlewares(b *testing.B) {
	uuid.EnableRandPool()
	defer func() {
		log.Printf("BENCHMARK: RequestLogging middleware called %d times", b.N)
	}()
	mainMiddleware := middleware.RegisterMiddlewares(
		middleware.RequestLogging,
	)
	router := http.NewServeMux()
	router.HandleFunc("GET /", dummyRequestHandler)
	mainMiddleware(router)
	for i := 0; i < b.N; i++ {
		_, _ = http.NewRequest("GET", "/", nil)
	}
}

func dummyRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(middleware.LoggerContextKey).(*slog.Logger)

	logger.Info("Logging from dummyRequestHandler", "status", "OK")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
