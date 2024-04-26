package model

type CreateBranchRequest struct {
	Name string `validate:"required" json:"name"`
}
