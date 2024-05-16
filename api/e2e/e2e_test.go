package e2e

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// TODO: write unit tests, this is just a rough guideline of what they should look like.
const baseURL = "http://localhost:8080"

func TestE2E(t *testing.T) {
	t.Run("HealthCheck", func(t *testing.T) {
		t.Run("should return healthy status", func(t *testing.T) {
			// Send a GET request to the health check endpoint
			resp, err := http.Get(baseURL + "/v1/health")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			defer resp.Body.Close()

			var healthResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&healthResp)
			assert.NoError(t, err)
			assert.Equal(t, "healthy", healthResp["status"])
		})
	})
}
