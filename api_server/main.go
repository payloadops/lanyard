// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Payload Ops API
 *
 * The Payload Ops API streamlines management of AI prompts, projects, organizations, teams, and users through conventional HTTP requests. This platform enables efficient automation and control of resources, providing robust tools for developers to manage settings, memberships, and activities seamlessly.
 *
 * API version: 1.0
 * Contact: info@payloadops.com
 */

package main

import (
	"log"
	"net/http"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID/go"
)

func main() {
	log.Printf("Server started")

	APIKeysAPIService := openapi.NewAPIKeysAPIService()
	APIKeysAPIController := openapi.NewAPIKeysAPIController(APIKeysAPIService)

	BranchesAPIService := openapi.NewBranchesAPIService()
	BranchesAPIController := openapi.NewBranchesAPIController(BranchesAPIService)

	CommitsAPIService := openapi.NewCommitsAPIService()
	CommitsAPIController := openapi.NewCommitsAPIController(CommitsAPIService)

	HealthCheckAPIService := openapi.NewHealthCheckAPIService()
	HealthCheckAPIController := openapi.NewHealthCheckAPIController(HealthCheckAPIService)

	OrganizationsAPIService := openapi.NewOrganizationsAPIService()
	OrganizationsAPIController := openapi.NewOrganizationsAPIController(OrganizationsAPIService)

	ProjectsAPIService := openapi.NewProjectsAPIService()
	ProjectsAPIController := openapi.NewProjectsAPIController(ProjectsAPIService)

	PromptsAPIService := openapi.NewPromptsAPIService()
	PromptsAPIController := openapi.NewPromptsAPIController(PromptsAPIService)

	TeamsAPIService := openapi.NewTeamsAPIService()
	TeamsAPIController := openapi.NewTeamsAPIController(TeamsAPIService)

	UsersAPIService := openapi.NewUsersAPIService()
	UsersAPIController := openapi.NewUsersAPIController(UsersAPIService)

	router := openapi.NewRouter(APIKeysAPIController, BranchesAPIController, CommitsAPIController, HealthCheckAPIController, OrganizationsAPIController, ProjectsAPIController, PromptsAPIController, TeamsAPIController, UsersAPIController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
