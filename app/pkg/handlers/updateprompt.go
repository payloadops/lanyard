package handlers

import (
	"errors"
	"net/http"
)

func validateUpdatePromptRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func UpdatePromptHandler(w http.ResponseWriter, r *http.Request) {

}
