package promptservice

type ListPromptsRequest struct {
	Name string `json:"name"`
}

type ListPromptsResponse struct {
	Name string `json:"name"`
}

func ListPrompts() []string {
	return []string{"test"}
}
