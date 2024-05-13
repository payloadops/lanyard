package userservicemodel

type GetUserResponse struct {
	Name string `validate:"required" json:"name"`
}
