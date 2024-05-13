package orghandler

import (
	"encoding/json"
	"net/http"
	"plato/app/pkg/model"
	"plato/app/pkg/service/apikey"
	"plato/app/pkg/util"
	keyservicemodel "plato/app_1/pkg/model/apikey/service"
	"strings"

	"github.com/go-chi/render"
)

func CreateOrgHandler(w http.ResponseWriter, r *http.Request) {
	validator := util.GetValidator()
	apiKeyService := apikey.NewService()

	var createApiKeyRequest keyservicemodel.CreateApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&createApiKeyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Struct(createApiKeyRequest); err != nil {
		render.Render(w, r, model.ErrorResponseRenderer(http.StatusBadRequest, err.Error()))
		return
	}

	urlSlices := strings.Split(r.URL.Path, "/")
	projectId := urlSlices[3]

	response, err := apiKeyService.Mint(
		r.Context(),
		projectId,
		createApiKeyRequest.Description,
		createApiKeyRequest.Scopes,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
