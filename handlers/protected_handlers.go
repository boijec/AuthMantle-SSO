package handlers

import (
	"authmantle-sso/middleware"
	"net/http"
)

// TODO remove this file?

func GetUserSettings(w http.ResponseWriter, r *http.Request) {
	tplCtx := r.Context().Value(middleware.TemplateContextKey).(*middleware.Templates)
	tplCtx.Render(w, "user_settings.html", Page{PageMeta: MetaData{PageTitle: "User Settings"}})
}
func GetAdminDashboardPage(w http.ResponseWriter, r *http.Request) {
	tplCtx := r.Context().Value(middleware.TemplateContextKey).(*middleware.Templates)
	tplCtx.Render(w, "admin_panel.html", Page{PageMeta: MetaData{PageTitle: "Admin Dashboard"}})
}
