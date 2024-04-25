package handlers

import (
	"encoding/json"
	"net/http"
	healthcheckservice "plato/app/pkg/service/health"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	response := healthcheckservice.CheckHealth()
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
