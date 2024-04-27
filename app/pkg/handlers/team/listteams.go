package teamhandler

import (
	"net/http"
	"plato/app/pkg/service/apikey"
	"strings"

	"github.com/go-chi/render"
)

func ListTeamsHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]

	response, err := apiKeyService.ListApiKeys(
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
