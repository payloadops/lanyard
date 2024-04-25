package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
	"strings"
)

func validateListPromptsRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func ListPromptsHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()

	if err := validateCreatePromptRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]

	response, err := promptService.ListPrompts(
		r.Context(),
		projectId,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
}
