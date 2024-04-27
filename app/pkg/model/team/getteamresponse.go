package teamservicemodel

type GetTeamResponse struct {
	Name string `validate:"required" json:"name"`
}
