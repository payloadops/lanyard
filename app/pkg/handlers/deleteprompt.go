package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
	"strings"
)

func validateDeletePromptRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func DeletePromptHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	promptService, _ := promptservice.NewService()

	if err := validateDeletePromptRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]
	promptId := urlSlices[5]

	response, err := promptService.DeletePrompt(
		r.Context(),
		projectId,
		promptId,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
