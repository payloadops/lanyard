package healthhandler

import (
	"github.com/go-chi/render"
	"net/http"
	"plato/app/pkg/service/health"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := health.NewService().CheckHealth()
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
