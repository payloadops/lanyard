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
	"time"
)

type Service struct {

	// Unique identifier for the service
	Id string `json:"id,omitempty"`

	// Name of the service
	Name string `json:"name,omitempty"`

	// A brief description of the service
	Description string `json:"description,omitempty"`

	// Timestamp when the service was created
	CreatedAt time.Time `json:"createdAt,omitempty"`

	// Timestamp when the service was last updated
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// AssertServiceRequired checks if the required fields are not zero-ed
func AssertServiceRequired(obj Service) error {
	return nil
}

// AssertServiceConstraints checks if the values respects the defined constraints
func AssertServiceConstraints(obj Service) error {
	return nil
}
