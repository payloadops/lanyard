package handlers

import (
	"errors"
	"net/http"
	"plato/app/pkg/service/promptversionservice"
)

func ValidateListPromptVersionsRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func ListPromptVersionsHandler(w http.ResponseWriter, r *http.Request) {
	promptversionservice.ListVersions()
}
