package projectservicemodel

type ListProjectsResponse struct {
	Name string `validate:"required" json:"name"`
}
