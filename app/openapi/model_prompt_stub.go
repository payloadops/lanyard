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

type PromptStub struct {

	// A K-sortable unique identifier (KSUID)
	Id string `json:"id,omitempty"`

	// Name of the prompt
	Name string `json:"name,omitempty"`

	// A brief description of the project
	Description string `json:"description,omitempty"`

	// A 180 character preview of the prompt
	Stub string `json:"stub,omitempty"`

	// A K-sortable unique identifier (KSUID)
	ProjectId string `json:"projectId,omitempty"`

	// Timestamp when the prompt was created
	CreatedAt time.Time `json:"createdAt,omitempty"`

	// Timestamp when the prompt was last updated
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// AssertPromptStubRequired checks if the required fields are not zero-ed
func AssertPromptStubRequired(obj PromptStub) error {
	return nil
}

// AssertPromptStubConstraints checks if the values respects the defined constraints
func AssertPromptStubConstraints(obj PromptStub) error {
	return nil
}
