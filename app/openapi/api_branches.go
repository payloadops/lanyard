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

// BranchesAPIController binds http requests to an api service and writes the service results to the http response
type BranchesAPIController struct {
	service BranchesAPIServicer
	errorHandler ErrorHandler
}

// BranchesAPIOption for how the controller is set up.
type BranchesAPIOption func(*BranchesAPIController)

// WithBranchesAPIErrorHandler inject ErrorHandler into controller
func WithBranchesAPIErrorHandler(h ErrorHandler) BranchesAPIOption {
	return func(c *BranchesAPIController) {
		c.errorHandler = h
	}
}

// NewBranchesAPIController creates a default api controller
func NewBranchesAPIController(s BranchesAPIServicer, opts ...BranchesAPIOption) Router {
	controller := &BranchesAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the BranchesAPIController
func (c *BranchesAPIController) Routes() Routes {
	return Routes{
		"CreatePromptBranch": Route{
			strings.ToUpper("Post"),
			"/v1/projects/{projectId}/prompts/{promptId}/branches",
			c.CreatePromptBranch,
		},
		"DeleteBranch": Route{
			strings.ToUpper("Delete"),
			"/v1/projects/{projectId}/prompts/{promptId}/branches/{branchName}",
			c.DeleteBranch,
		},
		"GetBranch": Route{
			strings.ToUpper("Get"),
			"/v1/projects/{projectId}/prompts/{promptId}/branches/{branchName}",
			c.GetBranch,
		},
		"ListPromptBranches": Route{
			strings.ToUpper("Get"),
			"/v1/projects/{projectId}/prompts/{promptId}/branches",
			c.ListPromptBranches,
		},
	}
}

// CreatePromptBranch - Create a new branch for a prompt
func (c *BranchesAPIController) CreatePromptBranch(w http.ResponseWriter, r *http.Request) {
	projectIdParam := chi.URLParam(r, "projectId")
	if projectIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"projectId"}, nil)
		return
	}
	promptIdParam := chi.URLParam(r, "promptId")
	if promptIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"promptId"}, nil)
		return
	}
	branchInputParam := BranchInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&branchInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertBranchInputRequired(branchInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertBranchInputConstraints(branchInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreatePromptBranch(r.Context(), projectIdParam, promptIdParam, branchInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// DeleteBranch - Delete a specific branch
func (c *BranchesAPIController) DeleteBranch(w http.ResponseWriter, r *http.Request) {
	projectIdParam := chi.URLParam(r, "projectId")
	if projectIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"projectId"}, nil)
		return
	}
	promptIdParam := chi.URLParam(r, "promptId")
	if promptIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"promptId"}, nil)
		return
	}
	branchNameParam := chi.URLParam(r, "branchName")
	if branchNameParam == "" {
		c.errorHandler(w, r, &RequiredError{"branchName"}, nil)
		return
	}
	result, err := c.service.DeleteBranch(r.Context(), projectIdParam, promptIdParam, branchNameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetBranch - Retrieve a specific branch
func (c *BranchesAPIController) GetBranch(w http.ResponseWriter, r *http.Request) {
	projectIdParam := chi.URLParam(r, "projectId")
	if projectIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"projectId"}, nil)
		return
	}
	promptIdParam := chi.URLParam(r, "promptId")
	if promptIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"promptId"}, nil)
		return
	}
	branchNameParam := chi.URLParam(r, "branchName")
	if branchNameParam == "" {
		c.errorHandler(w, r, &RequiredError{"branchName"}, nil)
		return
	}
	result, err := c.service.GetBranch(r.Context(), projectIdParam, promptIdParam, branchNameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ListPromptBranches - List all branches of a specific prompt
func (c *BranchesAPIController) ListPromptBranches(w http.ResponseWriter, r *http.Request) {
	projectIdParam := chi.URLParam(r, "projectId")
	if projectIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"projectId"}, nil)
		return
	}
	promptIdParam := chi.URLParam(r, "promptId")
	if promptIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"promptId"}, nil)
		return
	}
	result, err := c.service.ListPromptBranches(r.Context(), projectIdParam, promptIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
