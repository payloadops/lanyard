package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
	"strings"
)

func validateGetPromptRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func GetPromptHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()

	if err := validateGetPromptRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := promptService.GetPrompt(
		r.Context(),
		strings.Split(r.URL.Path, "/")[4],
		"main",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
