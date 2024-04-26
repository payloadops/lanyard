package keyservicemodel

type UpdateApiKeyRequest struct {
	Description string   `json:"description"`
	Scopes      []string `json:"scopes"`
}
