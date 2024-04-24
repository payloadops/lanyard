package model

type CreateApiKeyRequest struct {
	Description string   `json:"description"`
	Scopes      []string `json:"scopes"`
}
