package projectservicemodel

type GetProjectResponse struct {
	Name string `validate:"required" json:"name"`
}
