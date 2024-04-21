package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
)

type CreatePromptRequest struct {
	Name string `json:"name"`
}

// Validate checks the fields of UserData.
func (p *CreatePromptRequest) ValidateCreatePromptRequest() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// CreateUserHandler handles the user creation requests.
func CreatePromptHandler(w http.ResponseWriter, r *http.Request) {
	var promptData CreatePromptRequest
	if err := json.NewDecoder(r.Body).Decode(&promptData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := promptData.ValidateCreatePromptRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	promptservice.CreatePrompt()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(promptData)
}
