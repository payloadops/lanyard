package handlers

import (
	"errors"
	"net/http"
	"plato/app/pkg/service/promptversionservice"
)

type UpdateCurrentPromptVersionRequest struct {
	Name string `json:"name"`
}

func (p *UpdatePromptRequest) ValidateUpdateCurrentPromptVersionRequest() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func UpdateCurrentPromptVersionHandler(w http.ResponseWriter, r *http.Request) {
	promptversionservice.UpdateCurrentVersion()
}
