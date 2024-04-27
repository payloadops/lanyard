package teamhandler

import (
	"net/http"
	"plato/app/pkg/service/apikey"
	"strings"

	"github.com/go-chi/render"
)

func DeleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]
	apikey := urlSlices[5]

	err := apiKeyService.DeleteApiKey(
		r.Context(),
		projectId,
		apikey,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	// render.JSON(w, r, response)
}
