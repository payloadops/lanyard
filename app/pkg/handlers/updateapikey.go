package handlers

import (
	"encoding/json"
	"net/http"
	"plato/app/pkg/model"
	"plato/app/pkg/service/apikey"

	"github.com/go-chi/render"
)

func UpdateApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	apikey := r.URL.Path
	var updateApiKeyRequest model.CreateApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&updateApiKeyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := apiKeyService.UpdateApiKey(
		r.Context(),
		apikey,
		updateApiKeyRequest.Description,
		updateApiKeyRequest.Scopes,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	// render.JSON(w, r, response)
}
