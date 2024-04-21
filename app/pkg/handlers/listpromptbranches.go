package handlers

import (
	"errors"
	"net/http"
	"plato/app/pkg/service/promptbranchservice"
)

// Validate checks the fields of UserData.
func ValidateListPromptBranchesRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func ListPromptBranchesHandler(w http.ResponseWriter, r *http.Request) {
	promptbranchservice.ListBranches()
}
