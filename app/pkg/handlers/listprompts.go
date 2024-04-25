package handlers

import (
	"net/http"
	promptservice "plato/app/pkg/service/prompt"
	"strings"

	"github.com/go-chi/render"
)

func ListPromptsHandler(w http.ResponseWriter, r *http.Request) {
	promptService, _ := promptservice.NewService()

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]

	response, err := promptService.ListPrompts(
		r.Context(),
		projectId,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
