package projecthandler

import (
	"net/http"
	"plato/app/pkg/service/apikey"
	"strings"

	"github.com/go-chi/render"
)

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[4]
	apikey := urlSlices[6]

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
