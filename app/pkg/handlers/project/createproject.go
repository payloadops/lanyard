package projecthandler

import (
	"encoding/json"
	"net/http"
	projectservicemodel "plato/app/go/model/project"
	"plato/app/pkg/model"
	"plato/app/pkg/service/project"
	"plato/app/pkg/util"

	"github.com/go-chi/render"
)

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	validator := util.GetValidator()
	projectService, _ := project.NewService()

	var createProjectRequest projectservicemodel.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&createProjectRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Struct(createProjectRequest); err != nil {
		render.Render(w, r, model.ErrorResponseRenderer(http.StatusBadRequest, err.Error()))
		return
	}

	response, err := projectService.CreateProject(
		r.Context(),
		createProjectRequest,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
