package orgservicemodel

type DeleteOrgResponse struct {
	Name string `validate:"required" json:"name"`
}
