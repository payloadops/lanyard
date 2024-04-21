package promptbranchservice

type ListBranchesRequest struct {
	Name string `json:"name"`
}

type ListBranchesResponse struct {
	Name string `json:"name"`
}

func ListBranches() {

}
