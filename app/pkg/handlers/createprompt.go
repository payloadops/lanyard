package handlers

import (
	"encoding/json"
	"net/http"
	"plato/app/pkg/model"
	promptservice "plato/app/pkg/service/prompt"
	"plato/app/pkg/util"
	"strings"

	"github.com/go-chi/render"
)

func CreatePromptHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()
	validator := util.GetValidator()

	var createPromptRequest model.CreatePromptRequest
	if err := json.NewDecoder(r.Body).Decode(&createPromptRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Struct(createPromptRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]

	response, err := promptService.CreatePrompt(
		r.Context(),
		projectId,
		createPromptRequest,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
