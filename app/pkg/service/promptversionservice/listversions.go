package promptversionservice

type ListVersionsRequest struct {
	Name string `json:"name"`
}

type ListVersionsResponse struct {
	Name string `json:"name"`
}

func ListVersions() {

}
