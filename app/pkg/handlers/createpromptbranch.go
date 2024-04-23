package handlers

import (
	"errors"
	"net/http"
)

func validateCreatePromptBranchRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func CreatePromptBranchHandler(w http.ResponseWriter, r *http.Request) {

}
