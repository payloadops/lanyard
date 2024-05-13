package router

import (
	"net/http"
	healthhandler "plato/app_1/pkg/handlers/health"
	keyhandler "plato/app_1/pkg/handlers/key"
	orghandler "plato/app_1/pkg/handlers/org"
	projecthandler "plato/app_1/pkg/handlers/project"
	prompthandler "plato/app_1/pkg/handlers/prompt"
	teamhandler "plato/app_1/pkg/handlers/team"
	userhandler "plato/app_1/pkg/handlers/user"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const API_V1_PREFIX = "/api/v1"

func Init() {
	r := chi.NewRouter()
	setupMiddleware(r)
	setupRoutes(r)
	http.ListenAndServe(":8080", r)
}

func setupMiddleware(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(3 * time.Second))
	r.Use(middleware.AllowContentType("application/json", "charset=utf-8"))
	r.Use(middleware.Compress(5))

	r.Use(AuthMiddleware)
}

func setupRoutes(r *chi.Mux) {
	// Health check routes
	r.Get(API_V1_PREFIX+"/health", healthhandler.HealthCheckHandler)

	// API key management routes
	r.Post(API_V1_PREFIX+"/projects/{projectId}/api-keys", keyhandler.CreateApiKeyHandler)
	r.Get(API_V1_PREFIX+"/projects/{projectId}/api-keys", keyhandler.ListApiKeysHandler)
	r.Get(API_V1_PREFIX+"/projects/{projectId}/api-keys/{apiKey}", keyhandler.GetApiKeyHandler)
	r.Patch(API_V1_PREFIX+"/projects/{projectId}/api-keys/{apiKey}", keyhandler.UpdateApiKeyHandler)
	r.Delete(API_V1_PREFIX+"/projects/{projectId}/api-keys/{apiKey}", keyhandler.DeleteApiKeyHandler)

	// Prompt management routes
	r.Post(API_V1_PREFIX+"/projects/{projectId}/prompts", prompthandler.CreatePromptHandler)
	r.Get(API_V1_PREFIX+"/projects/{projectId}/prompts", prompthandler.ListPromptsHandler)

	r.Get(API_V1_PREFIX+"/projects/{projectId}/prompts/{promptId}", prompthandler.GetPromptHandler)
	r.Patch(API_V1_PREFIX+"/projects/{projectId}/prompts/{promptId}", prompthandler.UpdatePromptHandler)
	r.Delete(API_V1_PREFIX+"/projects/{projectId}/prompts/{promptId}", prompthandler.DeletePromptHandler)

	r.Get(API_V1_PREFIX+"/projects/{projectId}/prompts/{promptId}/versions", prompthandler.ListVersionsHandler)
	r.Put(API_V1_PREFIX+"/projects/{projectId}/prompts/{promptId}/versions", prompthandler.UpdateLiveVersionHandler)

	r.Post(API_V1_PREFIX+"/projects/{projectId}/prompts/{promptId}/branches", prompthandler.CreateBranchHandler)
	r.Get(API_V1_PREFIX+"/projects/{projectId}/prompts/{promptId}/branches", prompthandler.ListBranchesHandler)
	r.Delete(API_V1_PREFIX+"/projects/{projectId}/prompts/{promptId}/branches/{branchName}", prompthandler.DeleteBranchHandler)

	// Project management routes
	r.Get(API_V1_PREFIX+"/projects", projecthandler.ListProjectsHandler)
	r.Post(API_V1_PREFIX+"/projects", projecthandler.CreateProjectHandler)

	r.Get(API_V1_PREFIX+"/projects/{projectId}", projecthandler.GetProjectHandler)
	r.Patch(API_V1_PREFIX+"/projects/{projectId}", projecthandler.UpdateProjectHandler)
	r.Delete(API_V1_PREFIX+"/projects/{projectId}", projecthandler.DeleteProjectHandler)

	// Team management routes
	r.Get(API_V1_PREFIX+"/teams", teamhandler.ListTeamsHandler)
	r.Post(API_V1_PREFIX+"/teams", teamhandler.CreateTeamHandler)

	r.Get(API_V1_PREFIX+"/teams/{teamId}", teamhandler.GetTeamHandler)
	r.Patch(API_V1_PREFIX+"/teams/{teamId}", teamhandler.UpdateTeamHandler)
	r.Delete(API_V1_PREFIX+"/teams/{teamId}", teamhandler.DeleteTeamHandler)

	// User management routes
	r.Get(API_V1_PREFIX+"/users", userhandler.ListUsersHandler)
	r.Post(API_V1_PREFIX+"/users", userhandler.CreateUserHandler)

	r.Get(API_V1_PREFIX+"/users/{userId}", userhandler.GetUserHandler)
	r.Patch(API_V1_PREFIX+"/users/{userId}", userhandler.UpdateUserHandler)
	r.Delete(API_V1_PREFIX+"/users/{userId}", userhandler.DeleteUserHandler)

	// Org management routes
	r.Post(API_V1_PREFIX+"/orgs", orghandler.CreateOrgHandler)

	r.Get(API_V1_PREFIX+"/orgs/{orgId}", orghandler.GetOrgHandler)
	r.Patch(API_V1_PREFIX+"/orgs/{orgId}", orghandler.UpdateOrgHandler)
	r.Delete(API_V1_PREFIX+"/orgs/{orgId}", orghandler.DeleteOrgHandler)
}
