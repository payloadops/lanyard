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

func CreateBranchHandler(w http.ResponseWriter, r *http.Request) {
	validator := util.GetValidator()
	promptService, _ := promptservice.NewService()

	var createPromptBranchRequest model.CreateBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&createPromptBranchRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Struct(createPromptBranchRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]
	promptId := urlSlices[5]

	response, err := promptService.CreateBranch(
		r.Context(),
		projectId,
		promptId,
		createPromptBranchRequest,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
