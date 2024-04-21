package handlers

import (
	"errors"
	"net/http"
	"plato/app/pkg/service/promptversionservice"
)

func validateUpdateCurrentPromptVersionRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func UpdateCurrentPromptVersionHandler(w http.ResponseWriter, r *http.Request) {
	promptversionservice.UpdateCurrentVersion()
}
