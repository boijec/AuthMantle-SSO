package middleware

import (
	"net/http"
)

type ContextKey string
type Middleware func(http.Handler) http.Handler

func RegisterMiddlewares(m ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(m) - 1; i >= 0; i-- {
			j := m[i]
			next = j(next)
		}
		return next
	}
}
