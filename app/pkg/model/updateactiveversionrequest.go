package model

type UpdateActiveVersionRequest struct {
	Version string `validate:"required" json:"version,omitempty"`
	Branch  string `validate:"required" json:"branch,omitempty"`
}
