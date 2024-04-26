package keyservicemodel

type CreateApiKeyRequest struct {
	Description string   `validate:"required" json:"description"`
	Scopes      []string `validate:"required" json:"scopes"`
}
