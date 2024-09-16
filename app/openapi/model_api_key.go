// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Payload Ops API
 *
 * The Payload Ops API simplifies API key management for organizations by providing powerful tools to create, manage, and monitor API access securely. It allows teams to generate scoped API keys, configure rate limits, track usage, and integrate seamlessly with existing services.
 *
 * API version: 1.0
 * Contact: info@payloadops.com
 */

package openapi

import (
	"time"
)

type ApiKey struct {

	// Unique identifier for the API key
	Id string `json:"id,omitempty"`

	// The API key token
	Secret string `json:"secret,omitempty"`

	// List of roles granted by this API key
	Roles []string `json:"roles,omitempty"`

	// List of scopes or permissions granted by this API key
	Scopes []string `json:"scopes,omitempty"`

	// Number of monthly requests
	MonthlyRequestLimit int32 `json:"monthlyRequestLimit,omitempty"`

	// Rate limit configuration for this API key
	RateLimit RateLimit `json:"rateLimit,omitempty"`

	// The actor ID this API key is associated with
	ActorId string `json:"actorId,omitempty"`

	// The service ID this API key belongs to
	ServiceId string `json:"serviceId,omitempty"`

	// Timestamp when the API key was created
	CreatedAt time.Time `json:"createdAt,omitempty"`

	// Timestamp when the API key was last updated
	UpdatedAt time.Time `json:"updatedAt,omitempty"`

	// Name of the API key
	Name string `json:"name,omitempty"`

	// Billing information, including pricing tier and trial expiration date
	BillingInfo BillingInfo `json:"billingInfo,omitempty"`

	// Optional expiration date for the API key
	Expiry time.Time `json:"expiry,omitempty"`
}

// AssertApiKeyRequired checks if the required fields are not zero-ed
func AssertApiKeyRequired(obj ApiKey) error {
	if err := AssertRateLimitRequired(obj.RateLimit); err != nil {
		return err
	}
	if err := AssertBillingInfoRequired(obj.BillingInfo); err != nil {
		return err
	}
	return nil
}

// AssertApiKeyConstraints checks if the values respects the defined constraints
func AssertApiKeyConstraints(obj ApiKey) error {
	return nil
}
