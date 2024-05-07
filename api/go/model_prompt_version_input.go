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




type PromptVersionInput struct {

	// Content of the prompt at this version
	Content string `json:"content"`
}

// AssertPromptVersionInputRequired checks if the required fields are not zero-ed
func AssertPromptVersionInputRequired(obj PromptVersionInput) error {
	elements := map[string]interface{}{
		"content": obj.Content,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertPromptVersionInputConstraints checks if the values respects the defined constraints
func AssertPromptVersionInputConstraints(obj PromptVersionInput) error {
	return nil
}
