package handlers

import (
	"errors"
	"net/http"
	"plato/app/pkg/service/promptbranchservice"
)

func validateDeletePromptBranchRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func DeletePromptBranchHandler(w http.ResponseWriter, r *http.Request) {
	promptbranchservice.DeleteBranch()
}
