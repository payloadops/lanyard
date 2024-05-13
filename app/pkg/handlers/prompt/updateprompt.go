package prompthandler

import (
	"encoding/json"
	"net/http"
	promptservicemodel "plato/app/pkg/model/prompt/service"
	promptservice "plato/app/pkg/service/prompt"
	"plato/app/pkg/model"
	"plato/app/pkg/util"
	"strings"

	"github.com/go-chi/render"
)

func UpdatePromptHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()
	validator := util.GetValidator()

	var updatePromptRequest promptservicemodel.UpdatePromptRequest
	if err := json.NewDecoder(r.Body).Decode(&updatePromptRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Struct(updatePromptRequest); err != nil {
		render.Render(w, r, model.ErrorResponseRenderer(http.StatusBadRequest, err.Error()))
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[4]
	promptId := urlSlices[6]

	response, err := promptService.UpdatePrompt(
		r.Context(),
		projectId,
		promptId,
		updatePromptRequest,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
