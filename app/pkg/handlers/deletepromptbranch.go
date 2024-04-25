package handlers

import (
	"errors"
	"net/http"
)

func validateDeletePromptBranchRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func DeletePromptBranchHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
}
