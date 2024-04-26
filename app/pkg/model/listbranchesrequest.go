package model

type DeleteBranchResponse struct {
	Name string `validate:"required" json:"name"`
}
