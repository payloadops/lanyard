package projectservicemodel

type UpdateProjectRequest struct {
	TeamId      string `json:"team_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
