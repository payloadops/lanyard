package promptservicemodel

type DeletePromptResponse struct {
	PromptId  string `json:"prompt_id"`
	ProjectId string `json:"project_id"`
	DeletedAt string `json:"deleted_at"`
}
