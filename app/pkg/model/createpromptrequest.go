package model

type CreatePromptRequest struct {
	Name   string `json:"name"`
	Prompt string `json:"prompt"`
	Branch string `json:"branch"`
}
