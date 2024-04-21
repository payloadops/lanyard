package handlers

import (
	"errors"
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
)

type UpdatePromptRequest struct {
	Name string `json:"name"`
}

func (p *UpdatePromptRequest) ValidateUpdatePromptRequest() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func UpdatePromptHandler(w http.ResponseWriter, r *http.Request) {
	promptservice.UpdatePrompt()
}
