package handlers

import (
	"errors"
	"net/http"
	"plato/app/pkg/service/promptbranchservice"
)

func validateCreatePromptBranchRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func CreatePromptBranchHandler(w http.ResponseWriter, r *http.Request) {
	promptbranchservice.CreateBranch()
}
