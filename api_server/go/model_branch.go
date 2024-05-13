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



type Branch struct {

	// Unique identifier for the branch of the prompt
	Id string `json:"id,omitempty"`

	// Timestamp when this branch was created
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// AssertBranchRequired checks if the required fields are not zero-ed
func AssertBranchRequired(obj Branch) error {
	return nil
}

// AssertBranchConstraints checks if the values respects the defined constraints
func AssertBranchConstraints(obj Branch) error {
	return nil
}
