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

type ApiKeyInput struct {

	// List of roles granted by this API key
	Roles []string `json:"roles,omitempty"`

	// List of scopes or permissions granted by this API key
	Scopes []string `json:"scopes,omitempty"`

	// The actor ID this API key is associated with
	ActorExternalId string `json:"actorExternalId"`

	// Name of the API key
	Name string `json:"name"`

	// Optional expiration date for the API key
	Expiry time.Time `json:"expiry,omitempty"`
}

// AssertApiKeyInputRequired checks if the required fields are not zero-ed
func AssertApiKeyInputRequired(obj ApiKeyInput) error {
	elements := map[string]interface{}{
		"actorExternalId": obj.ActorExternalId,
		"name":            obj.Name,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertApiKeyInputConstraints checks if the values respects the defined constraints
func AssertApiKeyInputConstraints(obj ApiKeyInput) error {
	return nil
}
