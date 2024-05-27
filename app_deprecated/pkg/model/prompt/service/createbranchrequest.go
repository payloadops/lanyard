package promptservicemodel

type CreateBranchRequest struct {
	Name string `validate:"required" json:"name"`
}
