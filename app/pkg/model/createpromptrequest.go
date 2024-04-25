package model

type CreatePromptRequest struct {
	Name   string `validate:"required" json:"name"`
	Prompt string `validate:"required" json:"prompt"`
	Branch string `json:"branch"`
}
