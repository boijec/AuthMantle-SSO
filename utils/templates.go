package utils

import (
	"html/template"
	"net/http"
	"sync"
)

type Renderer struct {
	templates *template.Template
	lock      sync.Mutex
}

func (t *Renderer) Render(w http.ResponseWriter, name string, data interface{}) error {
	// TODO is this necessary?
	//t.lock.Lock()
	//defer t.lock.Unlock()
	return t.templates.ExecuteTemplate(w, name, data)
}

func InitializeTemplates() (Renderer, error) {
	return Renderer{
		// TODO export to env?
		templates: template.Must(template.ParseGlob("templates/*.html")),
		lock:      sync.Mutex{},
	}, nil
}
