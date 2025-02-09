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

// PricingTierAPIController binds http requests to an api service and writes the service results to the http response
type PricingTierAPIController struct {
	service      PricingTierAPIServicer
	errorHandler ErrorHandler
}

// PricingTierAPIOption for how the controller is set up.
type PricingTierAPIOption func(*PricingTierAPIController)

// WithPricingTierAPIErrorHandler inject ErrorHandler into controller
func WithPricingTierAPIErrorHandler(h ErrorHandler) PricingTierAPIOption {
	return func(c *PricingTierAPIController) {
		c.errorHandler = h
	}
}

// NewPricingTierAPIController creates a default api controller
func NewPricingTierAPIController(s PricingTierAPIServicer, opts ...PricingTierAPIOption) Router {
	controller := &PricingTierAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the PricingTierAPIController
func (c *PricingTierAPIController) Routes() Routes {
	return Routes{
		"ServicesServiceIdPricingTiersGet": Route{
			strings.ToUpper("Get"),
			"/v1/services/{serviceId}/pricing-tiers",
			c.ServicesServiceIdPricingTiersGet,
		},
		"ServicesServiceIdPricingTiersPost": Route{
			strings.ToUpper("Post"),
			"/v1/services/{serviceId}/pricing-tiers",
			c.ServicesServiceIdPricingTiersPost,
		},
		"ServicesServiceIdPricingTiersTierIdDelete": Route{
			strings.ToUpper("Delete"),
			"/v1/services/{serviceId}/pricing-tiers/{tierId}",
			c.ServicesServiceIdPricingTiersTierIdDelete,
		},
		"ServicesServiceIdPricingTiersTierIdGet": Route{
			strings.ToUpper("Get"),
			"/v1/services/{serviceId}/pricing-tiers/{tierId}",
			c.ServicesServiceIdPricingTiersTierIdGet,
		},
		"ServicesServiceIdPricingTiersTierIdPut": Route{
			strings.ToUpper("Put"),
			"/v1/services/{serviceId}/pricing-tiers/{tierId}",
			c.ServicesServiceIdPricingTiersTierIdPut,
		},
	}
}

// ServicesServiceIdPricingTiersGet - Retrieve the pricing tier for a service
func (c *PricingTierAPIController) ServicesServiceIdPricingTiersGet(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	result, err := c.service.ServicesServiceIdPricingTiersGet(r.Context(), serviceIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ServicesServiceIdPricingTiersPost - Assign a pricing tier to a service
func (c *PricingTierAPIController) ServicesServiceIdPricingTiersPost(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	pricingTierInputParam := PricingTierInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&pricingTierInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertPricingTierInputRequired(pricingTierInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertPricingTierInputConstraints(pricingTierInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ServicesServiceIdPricingTiersPost(r.Context(), serviceIdParam, pricingTierInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ServicesServiceIdPricingTiersTierIdDelete - Remove a pricing tier from a service
func (c *PricingTierAPIController) ServicesServiceIdPricingTiersTierIdDelete(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	tierIdParam := chi.URLParam(r, "tierId")
	if tierIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"tierId"}, nil)
		return
	}
	result, err := c.service.ServicesServiceIdPricingTiersTierIdDelete(r.Context(), serviceIdParam, tierIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ServicesServiceIdPricingTiersTierIdGet - Get the pricing tier for a service
func (c *PricingTierAPIController) ServicesServiceIdPricingTiersTierIdGet(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	tierIdParam := chi.URLParam(r, "tierId")
	if tierIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"tierId"}, nil)
		return
	}
	result, err := c.service.ServicesServiceIdPricingTiersTierIdGet(r.Context(), serviceIdParam, tierIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ServicesServiceIdPricingTiersTierIdPut - Update the pricing tier for a service
func (c *PricingTierAPIController) ServicesServiceIdPricingTiersTierIdPut(w http.ResponseWriter, r *http.Request) {
	serviceIdParam := chi.URLParam(r, "serviceId")
	if serviceIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"serviceId"}, nil)
		return
	}
	tierIdParam := chi.URLParam(r, "tierId")
	if tierIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"tierId"}, nil)
		return
	}
	pricingTierInputParam := PricingTierInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&pricingTierInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertPricingTierInputRequired(pricingTierInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertPricingTierInputConstraints(pricingTierInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ServicesServiceIdPricingTiersTierIdPut(r.Context(), serviceIdParam, tierIdParam, pricingTierInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
