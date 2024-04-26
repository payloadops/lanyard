package model

type ListBranchesResponse struct {
	Name string `validate:"required" json:"name"`
}
