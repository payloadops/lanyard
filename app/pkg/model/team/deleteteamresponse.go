package teamservicemodel

type DeleteTeamResponse struct {
	Name string `validate:"required" json:"name"`
}
