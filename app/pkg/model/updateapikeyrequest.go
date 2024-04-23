package model

type UpdateApiKeyRequest struct {
	Description string   `json:"description"`
	ProjectId   string   `json:"project_id"`
	Scopes      []string `json:"scopes"`
}
