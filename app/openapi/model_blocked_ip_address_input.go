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

// BlockedIpAddressInput - Information of blocked IP address and reason
type BlockedIpAddressInput struct {

	// IP address to be blocked
	IpAddress string `json:"ipAddress,omitempty"`

	// Reason why the IP address was blocked
	Reason string `json:"reason,omitempty"`
}

// AssertBlockedIpAddressInputRequired checks if the required fields are not zero-ed
func AssertBlockedIpAddressInputRequired(obj BlockedIpAddressInput) error {
	return nil
}

// AssertBlockedIpAddressInputConstraints checks if the values respects the defined constraints
func AssertBlockedIpAddressInputConstraints(obj BlockedIpAddressInput) error {
	return nil
}
