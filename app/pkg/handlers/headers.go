package handlers

import "net/http"

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
