package middleware

import (
	"authmantle-sso/data"
	"context"
	"net/http"
)

const SessionContextKey ContextKey = "user_session"

type AuthMiddleware struct {
	Db *data.DatabaseHandler
}

func (am *AuthMiddleware) EnsureSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("MANSESSION") // yeah, this is probably going to haunt me in the coming iterations
		if err != nil && cookie == nil {
			http.Redirect(w, r, "/error/401", http.StatusSeeOther) // TODO un-fuck this
			return
		}
		ctx := context.WithValue(r.Context(), SessionContextKey, cookie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
