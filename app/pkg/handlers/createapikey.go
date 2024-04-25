package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"plato/app/pkg/model"
	"plato/app/pkg/service/apikey"
	"strings"
)

func validateCreateApiKeyRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func CreateApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	err := validateCreateApiKeyRequest()

	apiKeyService := apikey.NewService()

	var createApiKeyRequest model.CreateApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&createApiKeyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateCreateApiKeyRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]

	response, err := apiKeyService.Mint(
		r.Context(),
		projectId,
		createApiKeyRequest.Description,
		createApiKeyRequest.Scopes,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
