package projectservicemodel

type CreateProjectResponse struct {
	Name        string `validate:"required" json:"name"`
	TeamId      string `validate:"required" json:"team_id"`
	Description string `validate:"required" json:"description"`
}
