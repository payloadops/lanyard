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
	ListBranchCommits(http.ResponseWriter, *http.Request)
}
// HealthCheckAPIRouter defines the required methods for binding the api requests to a responses for the HealthCheckAPI
// The HealthCheckAPIRouter implementation should parse necessary information from the http request,
// pass the data to a HealthCheckAPIServicer to perform the required actions, then write the service results to the http response.
type HealthCheckAPIRouter interface { 
	HealthCheck(http.ResponseWriter, *http.Request)
}
// OrganizationsAPIRouter defines the required methods for binding the api requests to a responses for the OrganizationsAPI
// The OrganizationsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a OrganizationsAPIServicer to perform the required actions, then write the service results to the http response.
type OrganizationsAPIRouter interface { 
	CreateOrganization(http.ResponseWriter, *http.Request)
	DeleteOrganization(http.ResponseWriter, *http.Request)
	GetOrganization(http.ResponseWriter, *http.Request)
	ListOrganizations(http.ResponseWriter, *http.Request)
	UpdateOrganization(http.ResponseWriter, *http.Request)
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
// TeamsAPIRouter defines the required methods for binding the api requests to a responses for the TeamsAPI
// The TeamsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a TeamsAPIServicer to perform the required actions, then write the service results to the http response.
type TeamsAPIRouter interface { 
	CreateTeam(http.ResponseWriter, *http.Request)
	DeleteTeam(http.ResponseWriter, *http.Request)
	GetTeam(http.ResponseWriter, *http.Request)
	ListTeams(http.ResponseWriter, *http.Request)
	UpdateTeam(http.ResponseWriter, *http.Request)
}
// UsersAPIRouter defines the required methods for binding the api requests to a responses for the UsersAPI
// The UsersAPIRouter implementation should parse necessary information from the http request,
// pass the data to a UsersAPIServicer to perform the required actions, then write the service results to the http response.
type UsersAPIRouter interface { 
	CreateUser(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
	ListUsers(http.ResponseWriter, *http.Request)
	UpdateUser(http.ResponseWriter, *http.Request)
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
	ListBranchCommits(context.Context, string, string, string) (ImplResponse, error)
}


// HealthCheckAPIServicer defines the api actions for the HealthCheckAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type HealthCheckAPIServicer interface { 
	HealthCheck(context.Context) (ImplResponse, error)
}


// OrganizationsAPIServicer defines the api actions for the OrganizationsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type OrganizationsAPIServicer interface { 
	CreateOrganization(context.Context, OrganizationInput) (ImplResponse, error)
	DeleteOrganization(context.Context, string) (ImplResponse, error)
	GetOrganization(context.Context, string) (ImplResponse, error)
	ListOrganizations(context.Context) (ImplResponse, error)
	UpdateOrganization(context.Context, string, OrganizationInput) (ImplResponse, error)
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


// TeamsAPIServicer defines the api actions for the TeamsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type TeamsAPIServicer interface { 
	CreateTeam(context.Context, string, TeamInput) (ImplResponse, error)
	DeleteTeam(context.Context, string, string) (ImplResponse, error)
	GetTeam(context.Context, string, string) (ImplResponse, error)
	ListTeams(context.Context, string) (ImplResponse, error)
	UpdateTeam(context.Context, string, string, TeamInput) (ImplResponse, error)
}


// UsersAPIServicer defines the api actions for the UsersAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type UsersAPIServicer interface { 
	CreateUser(context.Context, UserInput) (ImplResponse, error)
	DeleteUser(context.Context, string) (ImplResponse, error)
	GetUser(context.Context, string) (ImplResponse, error)
	ListUsers(context.Context) (ImplResponse, error)
	UpdateUser(context.Context, string, UserInput) (ImplResponse, error)
}
