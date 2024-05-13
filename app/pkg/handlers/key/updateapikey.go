package keyhandler

import (
	"encoding/json"
	"net/http"
	keyservicemodel "plato/app/openapi/model/apikey/service"
	"plato/app/pkg/service/apikey"
	"strings"

	"github.com/go-chi/render"
)

func UpdateApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	apiKeyService := apikey.NewService()

	var updateApiKeyRequest keyservicemodel.CreateApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&updateApiKeyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[4]
	apikey := urlSlices[6]

	err := apiKeyService.UpdateApiKey(
		r.Context(),
		projectId,
		apikey,
		updateApiKeyRequest.Description,
		updateApiKeyRequest.Scopes,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	// render.JSON(w, r, response)
}
