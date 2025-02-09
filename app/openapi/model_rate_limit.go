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

// RateLimit - Rate limit configuration for this API key
type RateLimit struct {

	// The name of the rate limit
	Name string `json:"name,omitempty"`

	// The number of allowed requests in the defined time window
	Limit int32 `json:"limit,omitempty"`

	// Time window for the rate limit, specified in ISO duration format (e.g., '1h', '30m')
	Duration string `json:"duration,omitempty"`
}

// AssertRateLimitRequired checks if the required fields are not zero-ed
func AssertRateLimitRequired(obj RateLimit) error {
	return nil
}

// AssertRateLimitConstraints checks if the values respects the defined constraints
func AssertRateLimitConstraints(obj RateLimit) error {
	return nil
}
