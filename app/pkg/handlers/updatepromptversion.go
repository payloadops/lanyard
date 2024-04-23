package handlers

import (
	"errors"
	"net/http"
)

func validateUpdateCurrentPromptVersionRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func UpdateCurrentPromptVersionHandler(w http.ResponseWriter, r *http.Request) {

}
