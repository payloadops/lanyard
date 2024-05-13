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
	"time"
)



type Team struct {

	// Unique identifier for the team
	Id string `json:"id,omitempty"`

	// Name of the team
	Name string `json:"name,omitempty"`

	// A brief description of the team
	Description string `json:"description,omitempty"`

	// The organization to which the team belongs
	OrgId string `json:"orgId,omitempty"`

	// Timestamp of team creation
	CreatedAt time.Time `json:"createdAt,omitempty"`

	// Timestamp of the last update to the team
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// AssertTeamRequired checks if the required fields are not zero-ed
func AssertTeamRequired(obj Team) error {
	return nil
}

// AssertTeamConstraints checks if the values respects the defined constraints
func AssertTeamConstraints(obj Team) error {
	return nil
}
