package orgservicemodel

type GetOrgResponse struct {
	Name string `validate:"required" json:"name"`
}
