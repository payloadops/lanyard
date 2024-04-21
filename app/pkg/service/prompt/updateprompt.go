package promptservice

type UpdatePromptRequest struct {
	Name string `json:"name"`
}

type UpdatePromptResponse struct {
	Name string `json:"name"`
}

func UpdatePrompt() string {
	return "test"
}
