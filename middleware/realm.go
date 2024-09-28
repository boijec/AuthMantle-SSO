package middleware

import (
	"authmantle-sso/data"
	"context"
	"log/slog"
	"net/http"
	"strings"
)

const RealmContextKey ContextKey = "realm"
const RealmIDContextKey ContextKey = "realm_id"

type RealmMiddleware struct {
	Db *data.DatabaseHandler
}

func (rm *RealmMiddleware) EnsureRealm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := ctx.Value(LoggerContextKey).(*slog.Logger)

		connection, err := rm.Db.Acquire(ctx)
		defer connection.Release()
		if err != nil {
			logger.Error("Internal error when acquiring database handle", "error", err)
			http.Redirect(w, r, "/error/500", http.StatusSeeOther)
			return
		}
		realmId := new(int)
		realmName := parseRealm(r.URL.Path)
		result := connection.QueryRow(ctx, "SELECT r.id FROM authmantledb.realm r WHERE r.name = $1", realmName)
		err = result.Scan(realmId)
		if err != nil || realmId == nil {
			logger.Error("Failed to find realm for request", "error", err)
			http.Redirect(w, r, "/error/404", http.StatusSeeOther)
			return
		}

		newCtx := context.WithValue(ctx, RealmIDContextKey, *realmId)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func parseRealm(rawUrl string) string {
	s := strings.Split(rawUrl, "/")[1]
	if len(s) > 50 {
		return ""
	}
	return s
}
