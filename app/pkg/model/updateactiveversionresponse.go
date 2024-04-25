package model

type UpdateActiveVersionResponse struct {
	Version string `validate:"required" json:"version"`
	Branch  string `validate:"required" json:"branch"`
}
