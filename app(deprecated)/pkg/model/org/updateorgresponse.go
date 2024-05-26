package orgservicemodel

type UpdateOrgResponse struct {
	Name        string `validate:"required" json:"name"`
	OrgId      string `validate:"required" json:"team_id"`
	Description string `validate:"required" json:"description"`
}
