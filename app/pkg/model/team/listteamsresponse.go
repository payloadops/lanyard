package teamservicemodel

type ListTeamsResponse struct {
	Name string `validate:"required" json:"name"`
}
