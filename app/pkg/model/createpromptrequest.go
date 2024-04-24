package model

type CreatePromptRequest struct {
	Prompt string `json:"prompt"`
	Branch string `json:"branch"`
}
