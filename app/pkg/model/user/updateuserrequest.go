package userservicemodel

type UpdateUserRequest struct {
	Name        string `validate:"required" json:"name"`
	TeamId      string `validate:"required" json:"team_id"`
	Description string `validate:"required" json:"description"`
}
