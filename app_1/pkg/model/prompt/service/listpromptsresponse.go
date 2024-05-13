package promptservicemodel

import dbdal "plato/app_1/pkg/dal/postgres"

type ListPromptsResponse struct {
	Prompts   *[]dbdal.Prompt `json:"prompts"`
	ProjectId string          `json:"project_id"`
}
