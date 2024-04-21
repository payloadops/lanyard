package handlers

import (
	"errors"
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
)

func ValidateGetPromptRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func GetPromptHandler(w http.ResponseWriter, r *http.Request) {
	promptservice.GetPrompt()
}
