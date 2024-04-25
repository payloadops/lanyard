package handlers

import (
	"encoding/json"
	"net/http"
	"plato/app/pkg/model"
	promptservice "plato/app/pkg/service/prompt"
	"strings"
)

func validatedUpdateCurrentPromptVersionRequest(w http.ResponseWriter, r *http.Request) (*model.UpdateActiveVersionRequest, error) {
	var updateActiveVersionRequest model.UpdateActiveVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&updateActiveVersionRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &updateActiveVersionRequest, nil
}

func setHeaders(w http.ResponseWriter) {

}

func UpdateCurrentPromptVersionHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	updateActiveVersionRequest, err := validatedUpdateCurrentPromptVersionRequest(w, r)
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
