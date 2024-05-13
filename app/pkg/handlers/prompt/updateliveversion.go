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

func UpdateLiveVersionHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()
	validator := util.GetValidator()

	var updateActiveVersionRequest promptservicemodel.UpdateActiveVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&updateActiveVersionRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Struct(updateActiveVersionRequest); err != nil {
		render.Render(w, r, model.ErrorResponseRenderer(http.StatusBadRequest, err.Error()))
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[4]
	promptId := urlSlices[6]

	response, err := promptService.UpdateActiveVersion(
		r.Context(),
		projectId,
		promptId,
		&updateActiveVersionRequest,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
