package model

type UpdateActiveVersionRequest struct {
	Version string `json:"version"`
	Branch  string `json:"branch"`
}
