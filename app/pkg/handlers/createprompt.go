package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"plato/app/pkg/model"
	promptservice "plato/app/pkg/service/prompt"
	"strings"
)

func validateCreatePromptRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func CreatePromptHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var createPromptRequest model.CreatePromptRequest
	promptService, _ := promptservice.NewService()
	if err := json.NewDecoder(r.Body).Decode(&createPromptRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateCreatePromptRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]

	response, err := promptService.CreatePrompt(
		r.Context(),
		projectId,
		createPromptRequest,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
