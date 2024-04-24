package model

import dbdal "plato/app/pkg/dal/postgres"

type ListPromptsResponse struct {
	Prompts   *[]dbdal.Prompt `json:"prompts"`
	ProjectId string          `json:"project_id"`
}
