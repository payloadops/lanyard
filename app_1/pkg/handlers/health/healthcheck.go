package healthhandler

import (
	"net/http"
	healthcheckservice "plato/app_1/pkg/service/health"

	"github.com/go-chi/render"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := healthcheckservice.CheckHealth()
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
