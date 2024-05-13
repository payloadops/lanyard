package prompthandler

import (
	"encoding/json"
	"net/http"
	promptservicemodel "plato/app/openapi/model/prompt/service"
	promptservice "plato/app/openapi/service/prompt"
	"plato/app/pkg/model"
	"plato/app/pkg/util"
	"strings"

	"github.com/go-chi/render"
)

func CreatePromptHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()
	validator := util.GetValidator()

	var createPromptRequest promptservicemodel.CreatePromptRequest
	if err := json.NewDecoder(r.Body).Decode(&createPromptRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Struct(createPromptRequest); err != nil {
		render.Render(w, r, model.ErrorResponseRenderer(http.StatusBadRequest, err.Error()))
		return
	}

	if len(createPromptRequest.Branch) == 0 {
		createPromptRequest.Branch = "main"
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[4]

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
