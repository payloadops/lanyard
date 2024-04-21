package handlers

import (
	"errors"
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
)

type DeletePromptRequest struct {
	Name string `json:"name"`
}

func (p *DeletePromptRequest) ValidateDeletePromptRequest() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func DeletePromptHandler(w http.ResponseWriter, r *http.Request) {
	promptservice.DeletePrompt()
}
