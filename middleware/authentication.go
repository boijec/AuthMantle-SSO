package middleware

import (
	"context"
	"net/http"
)

const SessionContextKey ContextKey = "user_session"

func EnsureSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("MANSESSION") // yeah, this is probably going to haunt me in the coming iterations
		if err != nil && cookie == nil {
			http.Redirect(w, r, "/v1/error/401", http.StatusSeeOther) // TODO un-fuck this
			return
		}
		ctx := context.WithValue(r.Context(), SessionContextKey, cookie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminLock(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("MANSESSION")
		if err != nil && cookie == nil {
			http.Redirect(w, r, "/adm_login/", http.StatusSeeOther) // smart, but dumb at the same time... TODO remove
			return
		}
		if cookie.Value != "adminBozo" {
			http.Redirect(w, r, "/adm_login/", http.StatusSeeOther) // TODO same here, remove
			return
		}
		ctx := context.WithValue(r.Context(), SessionContextKey, cookie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
