package userservicemodel

type ListUsersResponse struct {
	Name string `validate:"required" json:"name"`
}
