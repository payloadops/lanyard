package handlers

import (
	"errors"
	"net/http"
	"plato/app/pkg/service/promptbranchservice"
)

type DeletePromptBranchRequest struct {
	Name string `json:"name"`
}

// Validate checks the fields of UserData.
func (p *CreatePromptRequest) ValidateDeletePromptBranchRequest() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func DeletePromptBranchHandler(w http.ResponseWriter, r *http.Request) {
	promptbranchservice.DeleteBranch()
}
