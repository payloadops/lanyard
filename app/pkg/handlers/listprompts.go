package handlers

import (
	"errors"
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
)

func ValidateListPromptsRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func ListPromptsHandler(w http.ResponseWriter, r *http.Request) {
	promptservice.ListPrompts()
}
