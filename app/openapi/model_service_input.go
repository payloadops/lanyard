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

type ServiceInput struct {

	// Name of the service
	Name string `json:"name"`

	// A brief description of the service
	Description string `json:"description,omitempty"`
}

// AssertServiceInputRequired checks if the required fields are not zero-ed
func AssertServiceInputRequired(obj ServiceInput) error {
	elements := map[string]interface{}{
		"name": obj.Name,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertServiceInputConstraints checks if the values respects the defined constraints
func AssertServiceInputConstraints(obj ServiceInput) error {
	return nil
}
