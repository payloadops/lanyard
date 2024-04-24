package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"plato/app/pkg/service/apikey"
	"strings"
)

func validateListApiKeysRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func ListApiKeysHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	if err := validateListApiKeysRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]

	response, err := apiKeyService.ListApiKeys(
		r.Context(),
		projectId,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
