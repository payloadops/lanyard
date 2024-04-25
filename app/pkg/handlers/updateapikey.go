package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"plato/app/pkg/model"
	"plato/app/pkg/service/apikey"
)

func validateUpdateApiKeyRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func UpdateApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	apiKeyService := apikey.NewService()

	apikey := r.URL.Path
	var updateApiKeyRequest model.CreateApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&updateApiKeyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateUpdateApiKeyRequest(); err != nil {
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

	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(response)
}
