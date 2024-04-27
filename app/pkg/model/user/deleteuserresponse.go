package userservicemodel

type DeleteUserResponse struct {
	Name string `validate:"required" json:"name"`
}
