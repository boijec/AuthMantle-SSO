package middleware

import (
	"context"
	"html/template"
	"net/http"
)

type Templates struct {
	templates *template.Template
}

var Renderer Templates

//var Renderer = newTemplates()

func (t *Templates) Render(w http.ResponseWriter, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//func newTemplates() *Templates {
//	return &Templates{
//		// TODO export to env?
//		templates: template.Must(template.ParseGlob("templates/*.html")),
//	}
//}

const TemplateContextKey ContextKey = "renderer"

func RenderTemplateContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), TemplateContextKey, Renderer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
