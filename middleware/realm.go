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
type CacheWrapper struct {
	Realm string
}

// TODO call the db, once, ffs you're fucking it up again... use a timestamp to check if it needs to be refreshed.
// first caller will take the perf-hit.. TODO un-fuck this!
var cacheWrapper *CacheWrapper

func (rm *RealmMiddleware) EnsureRealm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		connection, err := rm.Db.Acquire(ctx)
		defer connection.Release()
		if err != nil {
			http.Redirect(w, r, "/error/500", http.StatusSeeOther)
			return
		}
		realmId := new(int32)
		result := connection.QueryRow(ctx, "SELECT r.id FROM authmantledb.realm r WHERE r.name = $1", parseRealm(r.URL.Path))
		err = result.Scan(realmId)
		if err != nil {
			http.Redirect(w, r, "/error/404", http.StatusSeeOther)
			return
		}

		newCtx := context.WithValue(ctx, RealmContextKey, r.PathValue("realm"))
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
