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

type PromptInput struct {

	// Name of the prompt
	Name string `json:"name"`

	// A brief description of the project
	Description string `json:"description,omitempty"`
}

// AssertPromptInputRequired checks if the required fields are not zero-ed
func AssertPromptInputRequired(obj PromptInput) error {
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

// AssertPromptInputConstraints checks if the values respects the defined constraints
func AssertPromptInputConstraints(obj PromptInput) error {
	return nil
}
