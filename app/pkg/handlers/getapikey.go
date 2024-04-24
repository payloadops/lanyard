package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"plato/app/pkg/service/apikey"
)

func validateGetApiKeyRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func GetApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	apikey := r.URL.Path

	if err := validateGetApiKeyRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := apiKeyService.GetApiKey(
		r.Context(),
		apikey,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
