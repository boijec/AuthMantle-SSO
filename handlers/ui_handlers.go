package handlers

import (
	"authmantle-sso/middleware"
	"net/http"
	"strconv"
)

// TODO yeah, make Page and MetaData more efficient this shouldn't be *that* hard

type Page struct {
	PageMeta   MetaData
	StatusCode int
	Error      string
}
type MetaData struct {
	PageTitle string
}

func GetLandingPage(w http.ResponseWriter, r *http.Request) {
	if s := r.URL.Path; s != "/" { // make sure that the shit does not effect other pages.
		http.Redirect(w, r, "/error/404", http.StatusSeeOther)
		return
	}
	tplCtx := r.Context().Value(middleware.TemplateContextKey).(*middleware.Templates)
	tplCtx.Render(w, "index.html", Page{PageMeta: MetaData{PageTitle: "Login"}})
}
func GetRegisterPage(w http.ResponseWriter, r *http.Request) {
	tplCtx := r.Context().Value(middleware.TemplateContextKey).(*middleware.Templates)
	tplCtx.Render(w, "register.html", Page{PageMeta: MetaData{PageTitle: "Login"}})
}
func GetAdminPage(w http.ResponseWriter, r *http.Request) {
	tplCtx := r.Context().Value(middleware.TemplateContextKey).(*middleware.Templates)
	tplCtx.Render(w, "admin_login.html", Page{PageMeta: MetaData{PageTitle: "Admin Login"}})
}
func ErrorRedirect(w http.ResponseWriter, r *http.Request) {
	status := parseStatusCode(r.PathValue("status"))
	tplCtx := r.Context().Value(middleware.TemplateContextKey).(*middleware.Templates)
	tplCtx.Render(w, "error.html", Page{PageMeta: MetaData{PageTitle: "Error"}, StatusCode: status})
}
func parseStatusCode(pathError string) int {
	if pathError == "" {
		return http.StatusInternalServerError
	}
	if len(pathError) > 4 {
		return http.StatusInternalServerError
	}
	status, err := strconv.Atoi(pathError)
	if err != nil {
		return http.StatusInternalServerError
	}

	return status
}
