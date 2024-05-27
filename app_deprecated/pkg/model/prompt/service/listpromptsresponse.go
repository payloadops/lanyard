package promptservicemodel

import dbdal "plato/app_deprecated/pkg/dal/postgres"

type ListPromptsResponse struct {
	Prompts   *[]dbdal.Prompt `json:"prompts"`
	ProjectId string          `json:"project_id"`
}
