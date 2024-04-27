package orgservicemodel

type ListOrgsResponse struct {
	Name string `validate:"required" json:"name"`
}
