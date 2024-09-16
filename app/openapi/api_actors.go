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

// ActorsAPIController binds http requests to an api service and writes the service results to the http response
type ActorsAPIController struct {
	service      ActorsAPIServicer
	errorHandler ErrorHandler
}

// ActorsAPIOption for how the controller is set up.
type ActorsAPIOption func(*ActorsAPIController)

// WithActorsAPIErrorHandler inject ErrorHandler into controller
func WithActorsAPIErrorHandler(h ErrorHandler) ActorsAPIOption {
	return func(c *ActorsAPIController) {
		c.errorHandler = h
	}
}

// NewActorsAPIController creates a default api controller
func NewActorsAPIController(s ActorsAPIServicer, opts ...ActorsAPIOption) Router {
	controller := &ActorsAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the ActorsAPIController
func (c *ActorsAPIController) Routes() Routes {
	return Routes{
		"ServicesServiceIdActorsActorIdDelete": Route{
			strings.ToUpper("Delete"),
			"/v1/services/{serviceId}/actors/{actorId}",
			c.ServicesServiceIdActorsActorIdDelete,
		},
		"ServicesServiceIdActorsActorIdGet": Route{
			strings.ToUpper("Get"),
			"/v1/services/{serviceId}/actors/{actorId}",
			c.ServicesServiceIdActorsActorIdGet,
		},
		"ServicesServiceIdActorsActorIdPut": Route{
			strings.ToUpper("Put"),
			"/v1/services/{serviceId}/actors/{actorId}",
			c.ServicesServiceIdActorsActorIdPut,
		},
		"ServicesServiceIdActorsGet": Route{
			strings.ToUpper("Get"),
			"/v1/services/{serviceId}/actors",
			c.ServicesServiceIdActorsGet,
		},
		"ServicesServiceIdActorsPost": Route{
			strings.ToUpper("Post"),
			"/v1/services/{serviceId}/actors",
			c.ServicesServiceIdActorsPost,
		},
	}
}

// ServicesServiceIdActorsActorIdDelete - Remove an actor from a service
func (c *ActorsAPIController) ServicesServiceIdActorsActorIdDelete(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	actorIdParam := chi.URLParam(r, "actorId")
	if actorIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"actorId"}, nil)
		return
	}
	result, err := c.service.ServicesServiceIdActorsActorIdDelete(r.Context(), serviceIdParam, actorIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ServicesServiceIdActorsActorIdGet - Get the actor
func (c *ActorsAPIController) ServicesServiceIdActorsActorIdGet(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	actorIdParam := chi.URLParam(r, "actorId")
	if actorIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"actorId"}, nil)
		return
	}
	result, err := c.service.ServicesServiceIdActorsActorIdGet(r.Context(), serviceIdParam, actorIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ServicesServiceIdActorsActorIdPut - Update an actor
func (c *ActorsAPIController) ServicesServiceIdActorsActorIdPut(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	actorIdParam := chi.URLParam(r, "actorId")
	if actorIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"actorId"}, nil)
		return
	}
	actorInputParam := ActorInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&actorInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertActorInputRequired(actorInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertActorInputConstraints(actorInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ServicesServiceIdActorsActorIdPut(r.Context(), serviceIdParam, actorIdParam, actorInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ServicesServiceIdActorsGet - Retrieve all actors associated with a service
func (c *ActorsAPIController) ServicesServiceIdActorsGet(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	result, err := c.service.ServicesServiceIdActorsGet(r.Context(), serviceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ServicesServiceIdActorsPost - Add an actor to a service
func (c *ActorsAPIController) ServicesServiceIdActorsPost(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	actorInputParam := ActorInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&actorInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertActorInputRequired(actorInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertActorInputConstraints(actorInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ServicesServiceIdActorsPost(r.Context(), serviceIdParam, actorInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
