package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"plato/app/pkg/service/apikey"
)

func validateListApiKeysRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func ListApiKeysHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	apikey := r.URL.Path

	if err := validateListApiKeysRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := apiKeyService.GetAPIKey(
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
