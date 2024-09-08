package middleware

import (
	"authmantle-sso/data"
	"context"
	"net/http"
)

const DbContextKey ContextKey = "db_connection"

func InjectDbContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := GetLogger(ctx)

		connection, err := data.GetFetcher().Acquire(ctx)
		defer func() {
			logger.DebugContext(ctx, "Releasing connection on account of finished request")
			connection.Release()
		}()
		if err != nil {
			logger.ErrorContext(ctx, "Failed to acquire database connection", "error", err)
			http.Error(w, "Failed to acquire database connection", http.StatusInternalServerError)
			return
		}
		ctx = context.WithValue(ctx, DbContextKey, connection)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
