package handlers

import (
	"errors"
	"net/http"
	"plato/app/pkg/service/apikey"
)

func validateDeleteApiKeyRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func DeleteApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	apikey := r.URL.Path
	if err := validateUpdateApiKeyRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := apiKeyService.DeleteAPIKey(
		r.Context(),
		apikey,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
