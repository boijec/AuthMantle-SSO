package middleware

import (
	"context"
	"github.com/google/uuid"
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

const RequestIDContextKey ContextKey = "user_session"

func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New()
		ctx = context.WithValue(ctx, RequestIDContextKey, id.String())
		r = r.WithContext(ctx)

		log.Printf("Incoming request [%s]: %s %s\n", id.String(), r.Method, r.URL.Path)
		s := time.Now()
		responseWrapper := &ResponseWrapper{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(responseWrapper, r)
		// TODO change to correlation id MW, when the request throws an error, this is dumb..
		log.Printf("Finished request [%s] with (%d) %s %s in %v\n", id.String(), responseWrapper.status, r.Method, r.URL.Path, time.Since(s))
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
