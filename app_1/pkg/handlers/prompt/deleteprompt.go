package prompthandler

import (
	"net/http"
	promptservice "plato/app_1/go/service/prompt"
	"strings"

	"github.com/go-chi/render"
)

func DeletePromptHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[4]
	promptId := urlSlices[6]

	response, err := promptService.DeletePrompt(
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
