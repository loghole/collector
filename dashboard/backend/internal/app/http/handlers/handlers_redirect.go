package handlers

import (
	"net/http"
)

func NewRedirectHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = path

		http.Redirect(w, r, r.URL.String(), http.StatusPermanentRedirect)
	}
}
