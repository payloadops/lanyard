// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Payload Ops API
 *
 * The Payload Ops API simplifies API key management for organizations by providing powerful tools to create, manage, and monitor API access securely. It allows teams to generate scoped API keys, configure rate limits, track usage, and integrate seamlessly with existing services.
 *
 * API version: 1.0
 * Contact: info@payloadops.com
 */

package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// APIKeysAPIController binds http requests to an api service and writes the service results to the http response
type APIKeysAPIController struct {
	service      APIKeysAPIServicer
	errorHandler ErrorHandler
}

// APIKeysAPIOption for how the controller is set up.
type APIKeysAPIOption func(*APIKeysAPIController)

// WithAPIKeysAPIErrorHandler inject ErrorHandler into controller
func WithAPIKeysAPIErrorHandler(h ErrorHandler) APIKeysAPIOption {
	return func(c *APIKeysAPIController) {
		c.errorHandler = h
	}
}

// NewAPIKeysAPIController creates a default api controller
func NewAPIKeysAPIController(s APIKeysAPIServicer, opts ...APIKeysAPIOption) Router {
	controller := &APIKeysAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the APIKeysAPIController
func (c *APIKeysAPIController) Routes() Routes {
	return Routes{
		"DeleteApiKey": Route{
			strings.ToUpper("Delete"),
			"/v1/services/{serviceId}/keys/{keyId}",
			c.DeleteApiKey,
		},
		"GenerateApiKey": Route{
			strings.ToUpper("Post"),
			"/v1/services/{serviceId}/keys",
			c.GenerateApiKey,
		},
		"GetApiKey": Route{
			strings.ToUpper("Get"),
			"/v1/services/{serviceId}/keys/{keyId}",
			c.GetApiKey,
		},
		"ListApiKeys": Route{
			strings.ToUpper("Get"),
			"/v1/services/{serviceId}/keys",
			c.ListApiKeys,
		},
		"UpdateApiKey": Route{
			strings.ToUpper("Put"),
			"/v1/services/{serviceId}/keys/{keyId}",
			c.UpdateApiKey,
		},
	}
}

// DeleteApiKey - Delete a specific API key
func (c *APIKeysAPIController) DeleteApiKey(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	keyIdParam := chi.URLParam(r, "keyId")
	if keyIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"keyId"}, nil)
		return
	}
	result, err := c.service.DeleteApiKey(r.Context(), serviceIdParam, keyIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GenerateApiKey - Generate a new API key with specific scopes for a service
func (c *APIKeysAPIController) GenerateApiKey(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	apiKeyInputParam := ApiKeyInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&apiKeyInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertApiKeyInputRequired(apiKeyInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertApiKeyInputConstraints(apiKeyInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.GenerateApiKey(r.Context(), serviceIdParam, apiKeyInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetApiKey - Retrieve a specific API key
func (c *APIKeysAPIController) GetApiKey(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	keyIdParam := chi.URLParam(r, "keyId")
	if keyIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"keyId"}, nil)
		return
	}
	result, err := c.service.GetApiKey(r.Context(), serviceIdParam, keyIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ListApiKeys - List all API keys for a service
func (c *APIKeysAPIController) ListApiKeys(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	result, err := c.service.ListApiKeys(r.Context(), serviceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateApiKey - Update an API key's scopes
func (c *APIKeysAPIController) UpdateApiKey(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	keyIdParam := chi.URLParam(r, "keyId")
	if keyIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"keyId"}, nil)
		return
	}
	apiKeyInputParam := ApiKeyInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&apiKeyInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertApiKeyInputRequired(apiKeyInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertApiKeyInputConstraints(apiKeyInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateApiKey(r.Context(), serviceIdParam, keyIdParam, apiKeyInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
