package promptservicemodel

type GetPromptResponse struct {
	Branch     string `json:"branch"`
	Prompt     string `json:"prompt"`
	Name       string `json:"name"`
	PromptId   string `json:"prompt_id"`
	Stub       string `json:"stub"`
	ProjectId  string `json:"project_id"`
	Version    string `json:"version"`
	ModifiedAt string `json:"modified_at"`
}
