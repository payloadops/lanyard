package prompthandler

import (
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
	"strings"

	"github.com/go-chi/render"
)

func DeleteBranchHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]
	promptId := urlSlices[5]
	branch := urlSlices[7]

	response, err := promptService.DeleteBranch(
		r.Context(),
		projectId,
		promptId,
		branch,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
