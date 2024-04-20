package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Healthy"}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
