package handlers

import (
	"net/http"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Set-Cookie", "MANSESSION=adminBozo; Path=/")
	http.Redirect(w, r, "/adm_console/", http.StatusSeeOther)
}
