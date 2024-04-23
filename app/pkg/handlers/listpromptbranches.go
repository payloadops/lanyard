package handlers

import (
	"errors"
	"net/http"
)

func validateListPromptBranchesRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func ListPromptBranchesHandler(w http.ResponseWriter, r *http.Request) {

}
