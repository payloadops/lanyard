package handlers

import (
	"errors"
	"net/http"
	"plato/app/pkg/service/promptbranchservice"
)

type CreatePromptBranchRequest struct {
	Name string `json:"name"`
}

// Validate checks the fields of UserData.
func (p *CreatePromptBranchRequest) ValidateCreatePromptBranchRequest() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func CreatePromptBranchHandler(w http.ResponseWriter, r *http.Request) {
	promptbranchservice.CreateBranch()
}
