package promptservice

type GetPromptRequest struct {
	Name string `json:"name"`
}

type GetPromptResponse struct {
	Name string `json:"name"`
}

func GetPrompt() string {
	return "test"
}
