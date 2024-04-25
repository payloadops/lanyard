package model

type UpdatePromptRequest struct {
	Branch string `json:"branch"`
	Prompt string `json:"prompt"`
	Name   string `json:"name"`
}
