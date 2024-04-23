package handlers

import (
	"errors"
	"net/http"
)

func validateGetPromptRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func GetPromptHandler(w http.ResponseWriter, r *http.Request) {

}
