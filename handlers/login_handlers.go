package handlers

import (
	"net/http"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Set-Cookie", "MANSESSION=adminBozo; Path=/") // note to self: you're an idiot...
	http.Redirect(w, r, "/adm_console/", http.StatusSeeOther)
}
