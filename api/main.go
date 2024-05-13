package main

import (
	"github.com/payloadops/plato/api/openapi"
	"github.com/payloadops/plato/api/service"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")

	APIKeysAPIService := service.NewAPIKeysAPIService()
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

	router := openapi.NewRouter(
		APIKeysAPIController,
		BranchesAPIController,
		CommitsAPIController,
		HealthCheckAPIController,
		OrganizationsAPIController,
		ProjectsAPIController,
		PromptsAPIController,
		TeamsAPIController,
		UsersAPIController,
	)

	log.Fatal(http.ListenAndServe(":8080", router))
}
