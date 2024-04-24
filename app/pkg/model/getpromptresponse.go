package model

type GetPromptResponse struct {
	Branch    string `json:"branch"`
	Prompt    string `json:"prompt"`
	PromptId  string `json:"prompt_id"`
	Stub      string `json:"stub"`
	ProjectId string `json:"project_id"`
	Version   string `json:"version"`
}
