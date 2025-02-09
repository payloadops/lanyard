// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Lanyard Ops API
 *
 * The Lanyard Ops API simplifies API key management for organizations by providing powerful tools to create, manage, and monitor API access securely. It allows teams to generate scoped API keys, configure rate limits, track usage, and integrate seamlessly with existing services.
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

// ServicesAPIController binds http requests to an api service and writes the service results to the http response
type ServicesAPIController struct {
	service      ServicesAPIServicer
	errorHandler ErrorHandler
}

// ServicesAPIOption for how the controller is set up.
type ServicesAPIOption func(*ServicesAPIController)

// WithServicesAPIErrorHandler inject ErrorHandler into controller
func WithServicesAPIErrorHandler(h ErrorHandler) ServicesAPIOption {
	return func(c *ServicesAPIController) {
		c.errorHandler = h
	}
}

// NewServicesAPIController creates a default api controller
func NewServicesAPIController(s ServicesAPIServicer, opts ...ServicesAPIOption) Router {
	controller := &ServicesAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the ServicesAPIController
func (c *ServicesAPIController) Routes() Routes {
	return Routes{
		"CreateService": Route{
			strings.ToUpper("Post"),
			"/v1/services",
			c.CreateService,
		},
		"DeleteService": Route{
			strings.ToUpper("Delete"),
			"/v1/services/{serviceId}",
			c.DeleteService,
		},
		"GetService": Route{
			strings.ToUpper("Get"),
			"/v1/services/{serviceId}",
			c.GetService,
		},
		"ListServices": Route{
			strings.ToUpper("Get"),
			"/v1/services",
			c.ListServices,
		},
		"UpdateService": Route{
			strings.ToUpper("Put"),
			"/v1/services/{serviceId}",
			c.UpdateService,
		},
	}
}

// CreateService - Create a new service
func (c *ServicesAPIController) CreateService(w http.ResponseWriter, r *http.Request) {
	serviceInputParam := ServiceInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&serviceInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertServiceInputRequired(serviceInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertServiceInputConstraints(serviceInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateService(r.Context(), serviceInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// DeleteService - Delete a service
func (c *ServicesAPIController) DeleteService(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	result, err := c.service.DeleteService(r.Context(), serviceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetService - Retrieve a service by ID
func (c *ServicesAPIController) GetService(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	result, err := c.service.GetService(r.Context(), serviceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ListServices - List all services
func (c *ServicesAPIController) ListServices(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.ListServices(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateService - Update a service
func (c *ServicesAPIController) UpdateService(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	serviceInputParam := ServiceInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&serviceInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertServiceInputRequired(serviceInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertServiceInputConstraints(serviceInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateService(r.Context(), serviceIdParam, serviceInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
