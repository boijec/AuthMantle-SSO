package utils

import (
	"authmantle-sso/data"
	"authmantle-sso/middleware"
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"sync"
)

type Renderer struct {
	templates *template.Template
	lock      sync.Mutex
}

type Page struct { // TODO move Countries from here
	PageMeta           MetaData
	RealmName          string
	EnableRegistration bool
	StatusCode         int
	Countries          []data.Country
	Error              string
}
type MetaData struct {
	PageTitle string
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func (t *Renderer) render(w http.ResponseWriter, name string, data interface{}) error {
	//t.lock.Lock()
	//defer t.lock.Unlock()
	return t.templates.ExecuteTemplate(w, name, data)
}

func InitializeTemplates() (Renderer, error) {
	tmp := template.New("master_renderer")
	tmp.Funcs(template.FuncMap{
		"getEnv": getEnv,
	})
	return Renderer{
		templates: template.Must(tmp.ParseGlob("web/templates/*.html")),
		lock:      sync.Mutex{},
	}, nil
}

func (t *Renderer) RenderErr(ctx context.Context, w http.ResponseWriter, name string, title string, error string) {
	logger := ctx.Value(middleware.LoggerContextKey).(*slog.Logger)
	err := t.render(w, name, Page{
		PageMeta:           MetaData{PageTitle: title},
		RealmName:          ctx.Value(middleware.RealmContextKey).(string),
		EnableRegistration: true,
		Error:              error,
	})
	if err != nil {
		http.Error(w, "Failed to process template", http.StatusInternalServerError)
		logger.ErrorContext(ctx, "Render failure", "error", err)
	}
}

func (t *Renderer) Render(ctx context.Context, w http.ResponseWriter, name string, title string) {
	logger := ctx.Value(middleware.LoggerContextKey).(*slog.Logger)
	renderData := Page{
		PageMeta:           MetaData{PageTitle: title},
		RealmName:          ctx.Value(middleware.RealmContextKey).(string),
		EnableRegistration: true,
	}
	err := t.render(w, name, renderData)
	if err != nil {
		http.Error(w, "Failed to process template", http.StatusInternalServerError)
		logger.ErrorContext(ctx, "Render failure", "error", err)
	}
}

func (t *Renderer) RenderWithData(ctx context.Context, w http.ResponseWriter, name string, requestData interface{}) {
	logger := ctx.Value(middleware.LoggerContextKey).(*slog.Logger)
	err := t.render(w, name, requestData)
	if err != nil {
		http.Error(w, "Failed to process template", http.StatusInternalServerError)
		logger.ErrorContext(ctx, "Render failure", "error", err)
	}
}
