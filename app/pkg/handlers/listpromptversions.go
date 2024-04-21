package handlers

import (
	"errors"
	"net/http"
	"plato/app/pkg/service/promptversionservice"
)

func validateListPromptVersionsRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func ListPromptVersionsHandler(w http.ResponseWriter, r *http.Request) {
	promptversionservice.ListVersions()
}
