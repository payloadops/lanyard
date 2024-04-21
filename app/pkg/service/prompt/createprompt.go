package promptservice

type CreatePromptRequest struct {
	Name string `json:"name"`
}

type CreatePromptResponse struct {
	Name string `json:"name"`
}

func CreatePrompt() string {
	return "test"
}
