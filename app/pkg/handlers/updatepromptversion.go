package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"plato/app/pkg/model"
	promptservice "plato/app/pkg/service/prompt"
	"strings"
)

func validateUpdateCurrentPromptVersionRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func UpdateCurrentPromptVersionHandler(w http.ResponseWriter, r *http.Request) {
	var updateActiveVersionRequest model.UpdateActiveVersionRequest
	promptService, _ := promptservice.NewService()
	if err := json.NewDecoder(r.Body).Decode(&updateActiveVersionRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]
	promptId := urlSlices[5]

	response, err := promptService.UpdateActiveVersion(
		r.Context(),
		projectId,
		promptId,
		updateActiveVersionRequest,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
