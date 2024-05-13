package prompthandler

import (
	"encoding/json"
	"net/http"
	"plato/app/pkg/model"
	"plato/app/pkg/util"
	promptservicemodel "plato/app_1/go/model/prompt/service"
	promptservice "plato/app_1/go/service/prompt"
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
