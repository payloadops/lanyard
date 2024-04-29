package prompthandler

import (
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
	"strings"

	"github.com/go-chi/render"
)

func GetPromptHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[4]
	promptId := urlSlices[6]
	query := r.URL.Query()
	branch := query.Get("branch")

	response, err := promptService.GetPrompt(
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
