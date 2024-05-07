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



type PromptDraft struct {

	// Unique identifier for the draft
	Id string `json:"id,omitempty"`

	// Name of the draft
	Name string `json:"name,omitempty"`

	// A brief description of the project
	Description string `json:"description,omitempty"`

	// Content of the draft
	Content string `json:"content,omitempty"`

	// Identifier of the prompt this draft belongs to
	PromptId string `json:"promptId,omitempty"`

	// Timestamp when the draft was created
	CreatedAt time.Time `json:"createdAt,omitempty"`

	// Timestamp when the draft was last updated
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// AssertPromptDraftRequired checks if the required fields are not zero-ed
func AssertPromptDraftRequired(obj PromptDraft) error {
	return nil
}

// AssertPromptDraftConstraints checks if the values respects the defined constraints
func AssertPromptDraftConstraints(obj PromptDraft) error {
	return nil
}
