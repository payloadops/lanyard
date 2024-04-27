package projectservicemodel

type DeleteProjectResponse struct {
	Name string `validate:"required" json:"name"`
}
