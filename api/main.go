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

	BranchesAPIService := service.NewBranchesAPIService()
	BranchesAPIController := openapi.NewBranchesAPIController(BranchesAPIService)

	CommitsAPIService := service.NewCommitsAPIService()
	CommitsAPIController := openapi.NewCommitsAPIController(CommitsAPIService)

	HealthCheckAPIService := service.NewHealthCheckAPIService()
	HealthCheckAPIController := openapi.NewHealthCheckAPIController(HealthCheckAPIService)

	OrganizationsAPIService := service.NewOrganizationsAPIService()
	OrganizationsAPIController := openapi.NewOrganizationsAPIController(OrganizationsAPIService)

	ProjectsAPIService := service.NewProjectsAPIService()
	ProjectsAPIController := openapi.NewProjectsAPIController(ProjectsAPIService)

	PromptsAPIService := service.NewPromptsAPIService()
	PromptsAPIController := openapi.NewPromptsAPIController(PromptsAPIService)

	TeamsAPIService := service.NewTeamsAPIService()
	TeamsAPIController := openapi.NewTeamsAPIController(TeamsAPIService)

	UsersAPIService := service.NewUsersAPIService()
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
