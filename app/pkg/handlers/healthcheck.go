package handlers

import (
	"encoding/json"
	"net/http"
	healthcheckservice "plato/app/pkg/service/health"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	healthcheckservice.CheckHealth()
	response := map[string]string{"message": "Healthy"}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
