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




type HealthCheck200Response struct {

	Status string `json:"status,omitempty"`

	Details HealthCheck200ResponseDetails `json:"details,omitempty"`
}

// AssertHealthCheck200ResponseRequired checks if the required fields are not zero-ed
func AssertHealthCheck200ResponseRequired(obj HealthCheck200Response) error {
	if err := AssertHealthCheck200ResponseDetailsRequired(obj.Details); err != nil {
		return err
	}
	return nil
}

// AssertHealthCheck200ResponseConstraints checks if the values respects the defined constraints
func AssertHealthCheck200ResponseConstraints(obj HealthCheck200Response) error {
	return nil
}
