// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Payload Ops API
 *
 * The Payload Ops API streamlines management of AI prompts, projects, organizations, teams, and users through conventional HTTP requests. This platform enables efficient automation and control of resources, providing robust tools for developers to manage settings, memberships, and activities seamlessly.
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

// DraftsAPIController binds http requests to an api service and writes the service results to the http response
type DraftsAPIController struct {
	service DraftsAPIServicer
	errorHandler ErrorHandler
}

// DraftsAPIOption for how the controller is set up.
type DraftsAPIOption func(*DraftsAPIController)

// WithDraftsAPIErrorHandler inject ErrorHandler into controller
func WithDraftsAPIErrorHandler(h ErrorHandler) DraftsAPIOption {
	return func(c *DraftsAPIController) {
		c.errorHandler = h
	}
}

// NewDraftsAPIController creates a default api controller
func NewDraftsAPIController(s DraftsAPIServicer, opts ...DraftsAPIOption) Router {
	controller := &DraftsAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the DraftsAPIController
func (c *DraftsAPIController) Routes() Routes {
	return Routes{
		"CreatePromptDraft": Route{
			strings.ToUpper("Post"),
			"/v1/prompts/{promptId}/drafts",
			c.CreatePromptDraft,
		},
		"ListPromptDrafts": Route{
			strings.ToUpper("Get"),
			"/v1/prompts/{promptId}/drafts",
			c.ListPromptDrafts,
		},
		"MergeDraftIntoVersion": Route{
			strings.ToUpper("Post"),
			"/v1/prompts/{promptId}/drafts/{draftId}/merge",
			c.MergeDraftIntoVersion,
		},
	}
}

// CreatePromptDraft - Create a new draft for a prompt
func (c *DraftsAPIController) CreatePromptDraft(w http.ResponseWriter, r *http.Request) {
	promptIdParam := chi.URLParam(r, "promptId")
	if promptIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"promptId"}, nil)
		return
	}
	promptDraftInputParam := PromptDraftInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&promptDraftInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertPromptDraftInputRequired(promptDraftInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertPromptDraftInputConstraints(promptDraftInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreatePromptDraft(r.Context(), promptIdParam, promptDraftInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ListPromptDrafts - List all drafts of a specific prompt
func (c *DraftsAPIController) ListPromptDrafts(w http.ResponseWriter, r *http.Request) {
	promptIdParam := chi.URLParam(r, "promptId")
	if promptIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"promptId"}, nil)
		return
	}
	result, err := c.service.ListPromptDrafts(r.Context(), promptIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// MergeDraftIntoVersion - Merge an approved draft into a new version
func (c *DraftsAPIController) MergeDraftIntoVersion(w http.ResponseWriter, r *http.Request) {
	promptIdParam := chi.URLParam(r, "promptId")
	if promptIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"promptId"}, nil)
		return
	}
	draftIdParam := chi.URLParam(r, "draftId")
	if draftIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"draftId"}, nil)
		return
	}
	result, err := c.service.MergeDraftIntoVersion(r.Context(), promptIdParam, draftIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
