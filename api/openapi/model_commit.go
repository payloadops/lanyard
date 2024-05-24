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



type Commit struct {

	// A K-sortable unique identifier (KSUID)
	Id string `json:"id,omitempty"`

	// Content of the prompt in this commit
	Content string `json:"content,omitempty"`

	// Message describing the changes in this commit
	Message string `json:"message,omitempty"`

	// SHA-256 checksum of the content in this commit
	Checksum string `json:"checksum,omitempty"`

	// A K-sortable unique identifier (KSUID)
	UserId string `json:"userId,omitempty"`

	// Timestamp when this commit was created
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// AssertCommitRequired checks if the required fields are not zero-ed
func AssertCommitRequired(obj Commit) error {
	return nil
}

// AssertCommitConstraints checks if the values respects the defined constraints
func AssertCommitConstraints(obj Commit) error {
	return nil
}
