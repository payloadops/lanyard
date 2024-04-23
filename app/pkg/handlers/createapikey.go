package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"plato/app/pkg/service/apikey"
)

type CreateApiKeyRequest struct {
	Description string   `json:"description"`
	ProjectId   string   `json:"project_id"`
	Scopes      []string `json:"scopes"`
}

func validateCreateApiKeyRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func CreateApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()
	var createApiKeyRequest CreateApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&createApiKeyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateCreateApiKeyRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := apiKeyService.Mint(
		r.Context(),
		createApiKeyRequest.ProjectId,
		createApiKeyRequest.Description,
		createApiKeyRequest.Scopes,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
