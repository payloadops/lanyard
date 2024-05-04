package router

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterSetup(t *testing.T) {
	r := NewRouter()

	// Define a test server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Define test cases for each route
	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/health"},
		{http.MethodPost, "/api/v1/projects/123/api-keys"},
		{http.MethodGet, "/api/v1/projects/123/api-keys"},
		{http.MethodGet, "/api/v1/projects/123/api-keys/abc"},
		{http.MethodPatch, "/api/v1/projects/123/api-keys/abc"},
		{http.MethodDelete, "/api/v1/projects/123/api-keys/abc"},
		{http.MethodPost, "/api/v1/projects/123/prompts"},
		{http.MethodGet, "/api/v1/projects/123/prompts"},
		{http.MethodGet, "/api/v1/projects/123/prompts/456"},
		{http.MethodPatch, "/api/v1/projects/123/prompts/456"},
		{http.MethodDelete, "/api/v1/projects/123/prompts/456"},
		{http.MethodGet, "/api/v1/projects/123/prompts/456/versions"},
		{http.MethodPut, "/api/v1/projects/123/prompts/456/versions"},
		{http.MethodPost, "/api/v1/projects/123/prompts/456/branches"},
		{http.MethodGet, "/api/v1/projects/123/prompts/456/branches"},
		{http.MethodDelete, "/api/v1/projects/123/prompts/456/branches/test"},
		{http.MethodGet, "/api/v1/projects"},
		{http.MethodPost, "/api/v1/projects"},
		{http.MethodGet, "/api/v1/projects/123"},
		{http.MethodPatch, "/api/v1/projects/123"},
		{http.MethodDelete, "/api/v1/projects/123"},
		{http.MethodGet, "/api/v1/teams"},
		{http.MethodPost, "/api/v1/teams"},
		{http.MethodGet, "/api/v1/teams/456"},
		{http.MethodPatch, "/api/v1/teams/456"},
		{http.MethodDelete, "/api/v1/teams/456"},
		{http.MethodGet, "/api/v1/users"},
		{http.MethodPost, "/api/v1/users"},
		{http.MethodGet, "/api/v1/users/789"},
		{http.MethodPatch, "/api/v1/users/789"},
		{http.MethodDelete, "/api/v1/users/789"},
		{http.MethodPost, "/api/v1/orgs"},
		{http.MethodGet, "/api/v1/orgs/abc"},
		{http.MethodPatch, "/api/v1/orgs/abc"},
		{http.MethodDelete, "/api/v1/orgs/abc"},
	}

	for _, test := range tests {
		req, err := http.NewRequest(test.method, ts.URL+test.path, nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		// TODO: make sure requests make their way to each respective handler.
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestRouterNotFound(t *testing.T) {
	r := NewRouter()

	// Define a test server
	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/foo", nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
