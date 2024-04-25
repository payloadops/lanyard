package handlers

import (
	"net/http"
	"plato/app/pkg/service/apikey"

	"github.com/go-chi/render"
)

func DeleteApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	apikey := r.URL.Path

	err := apiKeyService.DeleteApiKey(
		r.Context(),
		apikey,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	// render.JSON(w, r, response)
}
