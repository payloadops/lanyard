package projectservicemodel

type GetProjectResponse struct {
	ProjectId   string `json:"project_id"`
	OrgId       string `json:"orgId"`
	TeamId      string `json:"team_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
