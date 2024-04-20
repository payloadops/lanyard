package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type PromptData struct {
	Name string `json:"name"`
}

// Validate checks the fields of UserData.
func (p *PromptData) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// CreateUserHandler handles the user creation requests.
func CreatePromptHandler(w http.ResponseWriter, r *http.Request) {
	var promptData PromptData
	if err := json.NewDecoder(r.Body).Decode(&promptData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := promptData.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(promptData)
}

func ListPromptsHandler(w http.ResponseWriter, r *http.Request) {

}

func GetPromptHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdatePromptHandler(w http.ResponseWriter, r *http.Request) {

}

func DeletePromptHandler(w http.ResponseWriter, r *http.Request) {

}

func ListPromptVersionsHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdatePromptVersionHandler(w http.ResponseWriter, r *http.Request) {

}

func CreatePromptBranchHandler(w http.ResponseWriter, r *http.Request) {

}

func DeletePromptBranchHandler(w http.ResponseWriter, r *http.Request) {

}
