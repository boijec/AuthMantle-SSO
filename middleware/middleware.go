package middleware

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"
)

type ContextKey string
type Middleware func(http.Handler) http.Handler
type ResponseWrapper struct {
	http.ResponseWriter
	status int
}

func (rw *ResponseWrapper) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	rw.status = status
}

func RegisterMiddlewares(m ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(m) - 1; i >= 0; i-- {
			j := m[i]
			next = j(next)
		}
		return next
	}
}

func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()
		responseWrapper := &ResponseWrapper{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(responseWrapper, r)
		log.Printf("Request Finished in %v: (%d) %s %s\n", time.Since(s), responseWrapper.status, r.Method, r.URL.Path)
	})
}

type Templates struct {
	templates *template.Template
}

var Renderer = newTemplates()

func (t *Templates) Render(w http.ResponseWriter, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func newTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

const TemplateContextKey ContextKey = "renderer"
const SessionContextKey ContextKey = "user_session"

func RenderTemplateContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), TemplateContextKey, Renderer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func EnsureSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("MANSESSION")
		if err != nil && cookie == nil {
			http.Redirect(w, r, "/v1/error/401", http.StatusSeeOther)
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
			http.Redirect(w, r, "/adm_login/", http.StatusSeeOther)
			return
		}
		if cookie.Value != "adminBozo" {
			http.Redirect(w, r, "/adm_login/", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), SessionContextKey, cookie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
