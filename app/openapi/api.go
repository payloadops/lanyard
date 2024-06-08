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

import (
	"context"
	"net/http"
)

// APIKeysAPIRouter defines the required methods for binding the api requests to a responses for the APIKeysAPI
// The APIKeysAPIRouter implementation should parse necessary information from the http request,
// pass the data to a APIKeysAPIServicer to perform the required actions, then write the service results to the http response.
type APIKeysAPIRouter interface {
	DeleteApiKey(http.ResponseWriter, *http.Request)
	GenerateApiKey(http.ResponseWriter, *http.Request)
	GetApiKey(http.ResponseWriter, *http.Request)
	ListApiKeys(http.ResponseWriter, *http.Request)
	UpdateApiKey(http.ResponseWriter, *http.Request)
}

// BranchesAPIRouter defines the required methods for binding the api requests to a responses for the BranchesAPI
// The BranchesAPIRouter implementation should parse necessary information from the http request,
// pass the data to a BranchesAPIServicer to perform the required actions, then write the service results to the http response.
type BranchesAPIRouter interface {
	CreatePromptBranch(http.ResponseWriter, *http.Request)
	DeleteBranch(http.ResponseWriter, *http.Request)
	GetBranch(http.ResponseWriter, *http.Request)
	ListPromptBranches(http.ResponseWriter, *http.Request)
}

// CommitsAPIRouter defines the required methods for binding the api requests to a responses for the CommitsAPI
// The CommitsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a CommitsAPIServicer to perform the required actions, then write the service results to the http response.
type CommitsAPIRouter interface {
	CreateBranchCommit(http.ResponseWriter, *http.Request)
	GetBranchCommit(http.ResponseWriter, *http.Request)
	GetTemplateCommit(http.ResponseWriter, *http.Request)
	ListBranchCommits(http.ResponseWriter, *http.Request)
}

// HealthCheckAPIRouter defines the required methods for binding the api requests to a responses for the HealthCheckAPI
// The HealthCheckAPIRouter implementation should parse necessary information from the http request,
// pass the data to a HealthCheckAPIServicer to perform the required actions, then write the service results to the http response.
type HealthCheckAPIRouter interface {
	HealthCheck(http.ResponseWriter, *http.Request)
}

// ProjectsAPIRouter defines the required methods for binding the api requests to a responses for the ProjectsAPI
// The ProjectsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a ProjectsAPIServicer to perform the required actions, then write the service results to the http response.
type ProjectsAPIRouter interface {
	CreateProject(http.ResponseWriter, *http.Request)
	DeleteProject(http.ResponseWriter, *http.Request)
	GetProject(http.ResponseWriter, *http.Request)
	ListProjects(http.ResponseWriter, *http.Request)
	UpdateProject(http.ResponseWriter, *http.Request)
}

// PromptsAPIRouter defines the required methods for binding the api requests to a responses for the PromptsAPI
// The PromptsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a PromptsAPIServicer to perform the required actions, then write the service results to the http response.
type PromptsAPIRouter interface {
	CreatePrompt(http.ResponseWriter, *http.Request)
	DeletePrompt(http.ResponseWriter, *http.Request)
	GetPrompt(http.ResponseWriter, *http.Request)
	ListPrompts(http.ResponseWriter, *http.Request)
	UpdatePrompt(http.ResponseWriter, *http.Request)
}

// TestCasesAPIRouter defines the required methods for binding the api requests to a responses for the TestCasesAPI
// The TestCasesAPIRouter implementation should parse necessary information from the http request,
// pass the data to a TestCasesAPIServicer to perform the required actions, then write the service results to the http response.
type TestCasesAPIRouter interface {
	CreateTestCase(http.ResponseWriter, *http.Request)
	DeleteTestCase(http.ResponseWriter, *http.Request)
	GetTestCase(http.ResponseWriter, *http.Request)
	ListTestCases(http.ResponseWriter, *http.Request)
	UpdateTestCase(http.ResponseWriter, *http.Request)
}

// APIKeysAPIServicer defines the api actions for the APIKeysAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type APIKeysAPIServicer interface {
	DeleteApiKey(context.Context, string, string) (ImplResponse, error)
	GenerateApiKey(context.Context, string, ApiKeyInput) (ImplResponse, error)
	GetApiKey(context.Context, string, string) (ImplResponse, error)
	ListApiKeys(context.Context, string) (ImplResponse, error)
	UpdateApiKey(context.Context, string, string, ApiKeyInput) (ImplResponse, error)
}

// BranchesAPIServicer defines the api actions for the BranchesAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type BranchesAPIServicer interface {
	CreatePromptBranch(context.Context, string, string, BranchInput) (ImplResponse, error)
	DeleteBranch(context.Context, string, string, string) (ImplResponse, error)
	GetBranch(context.Context, string, string, string) (ImplResponse, error)
	ListPromptBranches(context.Context, string, string) (ImplResponse, error)
}

// CommitsAPIServicer defines the api actions for the CommitsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type CommitsAPIServicer interface {
	CreateBranchCommit(context.Context, string, string, string, CommitInput) (ImplResponse, error)
	GetBranchCommit(context.Context, string, string, string, string) (ImplResponse, error)
	GetTemplateCommit(context.Context, string, string, string) (ImplResponse, error)
	ListBranchCommits(context.Context, string, string, string) (ImplResponse, error)
}

// HealthCheckAPIServicer defines the api actions for the HealthCheckAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type HealthCheckAPIServicer interface {
	HealthCheck(context.Context) (ImplResponse, error)
}

// ProjectsAPIServicer defines the api actions for the ProjectsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ProjectsAPIServicer interface {
	CreateProject(context.Context, ProjectInput) (ImplResponse, error)
	DeleteProject(context.Context, string) (ImplResponse, error)
	GetProject(context.Context, string) (ImplResponse, error)
	ListProjects(context.Context) (ImplResponse, error)
	UpdateProject(context.Context, string, ProjectInput) (ImplResponse, error)
}

// PromptsAPIServicer defines the api actions for the PromptsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type PromptsAPIServicer interface {
	CreatePrompt(context.Context, string, PromptInput) (ImplResponse, error)
	DeletePrompt(context.Context, string, string) (ImplResponse, error)
	GetPrompt(context.Context, string, string) (ImplResponse, error)
	ListPrompts(context.Context, string) (ImplResponse, error)
	UpdatePrompt(context.Context, string, string, PromptInput) (ImplResponse, error)
}

// TestCasesAPIServicer defines the api actions for the TestCasesAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type TestCasesAPIServicer interface {
	CreateTestCase(context.Context, string, string, TestCaseInput) (ImplResponse, error)
	DeleteTestCase(context.Context, string, string, string) (ImplResponse, error)
	GetTestCase(context.Context, string, string, string) (ImplResponse, error)
	ListTestCases(context.Context, string, string) (ImplResponse, error)
	UpdateTestCase(context.Context, string, string, string, TestCaseInput) (ImplResponse, error)
}
