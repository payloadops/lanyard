package handlers

import (
	"errors"
	"net/http"
)

func validateDeletePromptRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func DeletePromptHandler(w http.ResponseWriter, r *http.Request) {

}
