package model

type CreateApiKeyRequest struct {
	Description string   `json:"description"`
	ProjectId   string   `json:"project_id"`
	Scopes      []string `json:"scopes"`
}
