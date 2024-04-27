package keyhandler

import (
	"net/http"
	"plato/app/pkg/service/apikey"
	"strings"

	"github.com/go-chi/render"
)

func GetApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	urlSlices := strings.Split(r.URL.Path, "/")
	apikey := urlSlices[6]

	response, err := apiKeyService.GetApiKey(
		r.Context(),
		apikey,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
