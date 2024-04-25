package handlers

import (
	"net/http"
	healthcheckservice "plato/app/pkg/service/health"

	"github.com/go-chi/render"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := healthcheckservice.CheckHealth()
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
