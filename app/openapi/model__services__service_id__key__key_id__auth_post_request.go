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

type ServicesServiceIdKeyKeyIdAuthPostRequest struct {

	// The API key provided by the client
	Secret string `json:"secret,omitempty"`

	//
	ActorExternalId string `json:"actorExternalId,omitempty"`

	// The API key provided by the client
	Roles string `json:"roles,omitempty"`

	// The API key provided by the client
	Scopes string `json:"scopes,omitempty"`
}

// AssertServicesServiceIdKeyKeyIdAuthPostRequestRequired checks if the required fields are not zero-ed
func AssertServicesServiceIdKeyKeyIdAuthPostRequestRequired(obj ServicesServiceIdKeyKeyIdAuthPostRequest) error {
	return nil
}

// AssertServicesServiceIdKeyKeyIdAuthPostRequestConstraints checks if the values respects the defined constraints
func AssertServicesServiceIdKeyKeyIdAuthPostRequestConstraints(obj ServicesServiceIdKeyKeyIdAuthPostRequest) error {
	return nil
}
