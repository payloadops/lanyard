package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plato/app/pkg/model"
	promptservice "plato/app/pkg/service/prompt"
	"strings"
)

func validatedUpdateCurrentPromptVersionRequest(w http.ResponseWriter, r *http.Request) (*model.UpdateActiveVersionRequest, error) {
	var updateActiveVersionRequest model.UpdateActiveVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&updateActiveVersionRequest); err != nil {
		return nil, err
	}
	if len(updateActiveVersionRequest.Branch) == 0 {
		return nil, fmt.Errorf("branch is in incorrect format")
	}
	if len(updateActiveVersionRequest.Version) == 0 {
		return nil, fmt.Errorf("version is in incorrect format")
	}
	return &updateActiveVersionRequest, nil
}

func UpdateCurrentPromptVersionHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	updateActiveVersionRequest, err := validatedUpdateCurrentPromptVersionRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	promptService, _ := promptservice.NewService()

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
