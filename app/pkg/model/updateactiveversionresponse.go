package model

type UpdateActiveVersionResponse struct {
	Version string `json:"version"`
	Branch  string `json:"branch"`
}
