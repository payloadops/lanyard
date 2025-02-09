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

type AuthApiKey200ResponseRateLimit struct {
	Reset int32 `json:"reset,omitempty"`

	Limit int32 `json:"limit,omitempty"`

	Remaining int32 `json:"remaining,omitempty"`
}

// AssertAuthApiKey200ResponseRateLimitRequired checks if the required fields are not zero-ed
func AssertAuthApiKey200ResponseRateLimitRequired(obj AuthApiKey200ResponseRateLimit) error {
	return nil
}

// AssertAuthApiKey200ResponseRateLimitConstraints checks if the values respects the defined constraints
func AssertAuthApiKey200ResponseRateLimitConstraints(obj AuthApiKey200ResponseRateLimit) error {
	return nil
}
