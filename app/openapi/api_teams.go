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

// TeamsAPIController binds http requests to an api service and writes the service results to the http response
type TeamsAPIController struct {
	service      TeamsAPIServicer
	errorHandler ErrorHandler
}

// TeamsAPIOption for how the controller is set up.
type TeamsAPIOption func(*TeamsAPIController)

// WithTeamsAPIErrorHandler inject ErrorHandler into controller
func WithTeamsAPIErrorHandler(h ErrorHandler) TeamsAPIOption {
	return func(c *TeamsAPIController) {
		c.errorHandler = h
	}
}

// NewTeamsAPIController creates a default api controller
func NewTeamsAPIController(s TeamsAPIServicer, opts ...TeamsAPIOption) Router {
	controller := &TeamsAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the TeamsAPIController
func (c *TeamsAPIController) Routes() Routes {
	return Routes{
		"CreateTeam": Route{
			strings.ToUpper("Post"),
			"/v1/organizations/{orgId}/teams",
			c.CreateTeam,
		},
		"DeleteTeam": Route{
			strings.ToUpper("Delete"),
			"/v1/organizations/{orgId}/teams/{teamId}",
			c.DeleteTeam,
		},
		"GetTeam": Route{
			strings.ToUpper("Get"),
			"/v1/organizations/{orgId}/teams/{teamId}",
			c.GetTeam,
		},
		"ListTeams": Route{
			strings.ToUpper("Get"),
			"/v1/organizations/{orgId}/teams",
			c.ListTeams,
		},
		"UpdateTeam": Route{
			strings.ToUpper("Put"),
			"/v1/organizations/{orgId}/teams/{teamId}",
			c.UpdateTeam,
		},
	}
}

// CreateTeam - Create a new team for an organization
func (c *TeamsAPIController) CreateTeam(w http.ResponseWriter, r *http.Request) {
	orgIdParam := chi.URLParam(r, "orgId")
	if orgIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"orgId"}, nil)
		return
	}
	teamInputParam := TeamInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&teamInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTeamInputRequired(teamInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertTeamInputConstraints(teamInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateTeam(r.Context(), orgIdParam, teamInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// DeleteTeam - Delete a specific team
func (c *TeamsAPIController) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	orgIdParam := chi.URLParam(r, "orgId")
	if orgIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"orgId"}, nil)
		return
	}
	teamIdParam := chi.URLParam(r, "teamId")
	if teamIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"teamId"}, nil)
		return
	}
	result, err := c.service.DeleteTeam(r.Context(), orgIdParam, teamIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetTeam - Get a specific team by ID
func (c *TeamsAPIController) GetTeam(w http.ResponseWriter, r *http.Request) {
	orgIdParam := chi.URLParam(r, "orgId")
	if orgIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"orgId"}, nil)
		return
	}
	teamIdParam := chi.URLParam(r, "teamId")
	if teamIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"teamId"}, nil)
		return
	}
	result, err := c.service.GetTeam(r.Context(), orgIdParam, teamIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ListTeams - List all teams in an organization
func (c *TeamsAPIController) ListTeams(w http.ResponseWriter, r *http.Request) {
	orgIdParam := chi.URLParam(r, "orgId")
	if orgIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"orgId"}, nil)
		return
	}
	result, err := c.service.ListTeams(r.Context(), orgIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateTeam - Update a specific team
func (c *TeamsAPIController) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	orgIdParam := chi.URLParam(r, "orgId")
	if orgIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"orgId"}, nil)
		return
	}
	teamIdParam := chi.URLParam(r, "teamId")
	if teamIdParam == "" {
		c.errorHandler(w, r, &RequiredError{"teamId"}, nil)
		return
	}
	teamInputParam := TeamInput{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&teamInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTeamInputRequired(teamInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertTeamInputConstraints(teamInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateTeam(r.Context(), orgIdParam, teamIdParam, teamInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
