package projectservicemodel

import dbdal "plato/app_1/pkg/dal/postgres"

type ListProjectsResponse struct {
	OrgId    string           `json:"org_id"`
	TeamId   string           `json:"team_id"`
	Projects *[]dbdal.Project `json:"projects"`
}
