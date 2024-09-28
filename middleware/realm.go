package middleware

import (
	"authmantle-sso/data"
	"context"
	"net/http"
	"strings"
)

const RealmContextKey ContextKey = "realm"

type RealmMiddleware struct {
	Db *data.DatabaseHandler
}

func (rm *RealmMiddleware) EnsureRealm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		connection, err := rm.Db.Acquire(ctx)
		defer connection.Release()
		if err != nil {
			http.Redirect(w, r, "/error/500", http.StatusSeeOther)
			return
		}
		realmId := new(int)
		realmName := parseRealm(r.URL.Path)
		result := connection.QueryRow(ctx, "SELECT count(*) FROM authmantledb.realm r WHERE r.name = $1", realmName)
		err = result.Scan(realmId)
		if err != nil || *realmId != 1 {
			http.Redirect(w, r, "/error/404", http.StatusSeeOther)
			return
		}

		newCtx := context.WithValue(ctx, RealmContextKey, realmName)
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
