package e2e

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// TODO: write unit tests, this is just a rough guideline of what they should look like.
const baseURL = "http://localhost:8080"

func TestE2E(t *testing.T) {
	t.Run("Organizations", func(t *testing.T) {
		t.Run("should create and retrieve an organization", func(t *testing.T) {
			// Create an organization
			orgPayload := map[string]interface{}{
				"name":        "Test Organization",
				"description": "An organization for testing",
			}
			orgBody, _ := json.Marshal(orgPayload)
			resp, err := http.Post(baseURL+"/organizations", "application/json", bytes.NewBuffer(orgBody))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			defer resp.Body.Close()
			var createdOrg map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&createdOrg)
			assert.NoError(t, err)

			orgID := createdOrg["id"].(string)
			assert.NotEmpty(t, orgID)

			// Retrieve the created organization
			resp, err = http.Get(baseURL + "/organizations/" + orgID)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			defer resp.Body.Close()
			var retrievedOrg map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&retrievedOrg)
			assert.NoError(t, err)
			assert.Equal(t, orgID, retrievedOrg["id"])
			assert.Equal(t, orgPayload["name"], retrievedOrg["name"])
			assert.Equal(t, orgPayload["description"], retrievedOrg["description"])
		})
	})

	t.Run("Projects", func(t *testing.T) {
		var orgID string

		// Setup: Create an organization for the project
		orgPayload := map[string]interface{}{
			"name":        "Test Organization",
			"description": "An organization for testing",
		}
		orgBody, _ := json.Marshal(orgPayload)
		resp, err := http.Post(baseURL+"/organizations", "application/json", bytes.NewBuffer(orgBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()
		var createdOrg map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&createdOrg)
		assert.NoError(t, err)

		orgID = createdOrg["id"].(string)
		assert.NotEmpty(t, orgID)

		t.Run("should create and retrieve a project", func(t *testing.T) {
			// Create a project
			projectPayload := map[string]interface{}{
				"name":        "Test Project",
				"description": "A project for testing",
				"orgId":       orgID,
			}
			projectBody, _ := json.Marshal(projectPayload)
			resp, err := http.Post(baseURL+"/projects", "application/json", bytes.NewBuffer(projectBody))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			defer resp.Body.Close()
			var createdProject map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&createdProject)
			assert.NoError(t, err)

			projectID := createdProject["id"].(string)
			assert.NotEmpty(t, projectID)

			// Retrieve the created project
			resp, err = http.Get(baseURL + "/projects/" + projectID)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			defer resp.Body.Close()
			var retrievedProject map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&retrievedProject)
			assert.NoError(t, err)
			assert.Equal(t, projectID, retrievedProject["id"])
			assert.Equal(t, projectPayload["name"], retrievedProject["name"])
			assert.Equal(t, projectPayload["description"], retrievedProject["description"])
			assert.Equal(t, orgID, retrievedProject["orgId"])
		})
	})
}
