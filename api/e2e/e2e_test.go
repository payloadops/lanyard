package e2e

import (
	"bytes"
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// Config holds the configuration for the e2e tests.
type Config struct {
	BaseURL string `envconfig:"BASE_URL" default:"http://localhost:8080"`
}

// AuthSchemeType is a string enum for authentication schemes.
type AuthSchemeType string

const (
	AuthSchemeNone   AuthSchemeType = "none"
	AuthSchemeBearer AuthSchemeType = "bearer"
	AuthSchemeApiKey AuthSchemeType = "apiKey"
)

// TestConfig defines a test case.
type TestConfig struct {
	Name           string
	Method         string
	Endpoint       string
	Body           map[string]interface{}
	ExpectedStatus int
	ExpectedBody   map[string]interface{}
	AuthScheme     AuthSchemeType
	AuthToken      string
	ApiKey         string
}

// runTest executes a single test case.
func runTest(t *testing.T, baseURL string, tc TestConfig) {
	var req *http.Request
	var err error

	fullURL := baseURL + tc.Endpoint
	if tc.Body != nil {
		body, _ := json.Marshal(tc.Body)
		req, err = http.NewRequest(tc.Method, fullURL, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(tc.Method, fullURL, nil)
	}

	assert.NoError(t, err)

	// Set the appropriate authentication header based on the scheme
	switch tc.AuthScheme {
	case AuthSchemeBearer:
		if tc.AuthToken != "" {
			req.Header.Set("Authorization", "Bearer "+tc.AuthToken)
		}
	case AuthSchemeApiKey:
		if tc.ApiKey != "" {
			req.Header.Set("X-API-KEY", tc.ApiKey)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)

	defer resp.Body.Close()

	if tc.ExpectedBody != nil {
		var actualBody map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&actualBody)
		assert.NoError(t, err)
		assert.Equal(t, tc.ExpectedBody, actualBody)
	}
}

// TestE2E runs the end-to-end tests defined in the TestConfig array.
func TestE2E(t *testing.T) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	assert.NoError(t, err)

	tests := []TestConfig{
		{
			Name:           "HealthCheck should return healthy status",
			Method:         http.MethodGet,
			Endpoint:       "/v1/health",
			ExpectedStatus: http.StatusOK,
			ExpectedBody: map[string]interface{}{
				"status": "healthy",
			},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			runTest(t, cfg.BaseURL, tc)
		})
	}
}
