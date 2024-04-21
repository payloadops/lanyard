package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
)

func validateCreatePromptRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func CreatePromptHandler(w http.ResponseWriter, r *http.Request) {
	var createPromptRequest promptservice.CreatePromptRequest
	if err := json.NewDecoder(r.Body).Decode(&createPromptRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateCreatePromptRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	promptservice.CreatePrompt()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createPromptRequest)
}
