package prompthandler

import (
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
	"strings"

	"github.com/go-chi/render"
)

func ListBranchesHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]
	promptId := urlSlices[5]

	response, err := promptService.ListBranches(
		r.Context(),
		projectId,
		promptId,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
