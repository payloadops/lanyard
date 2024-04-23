package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

func validateCreatePromptRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func CreatePromptHandler(w http.ResponseWriter, r *http.Request) {
	// promptService, err := promptservice.NewService()
	if err := json.NewDecoder(r.Body).Decode(""); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateCreatePromptRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// response := promptService.CreatePrompt()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("response")
}
