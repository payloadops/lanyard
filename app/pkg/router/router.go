package router

import (
	healthhandler "plato/app/pkg/handlers/health"
	keyhandler "plato/app/pkg/handlers/key"
	orghandler "plato/app/pkg/handlers/org"
	projecthandler "plato/app/pkg/handlers/project"
	prompthandler "plato/app/pkg/handlers/prompt"
	teamhandler "plato/app/pkg/handlers/team"
	userhandler "plato/app/pkg/handlers/user"
	authmiddleware "plato/app/pkg/middleware/auth"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ApiV1Prefix defines the API prefix for version 1.
const ApiV1Prefix = "/api/v1"

// ApiTimeout defines the timeout for API requests.
const ApiTimeout = 3 * time.Second

// NewRouter creates and returns a new chi.Mux router with configured middleware and routes.
func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	setupMiddleware(r)
	setupRoutes(r)
	return r
}

// setupMiddleware configures and adds middleware to the router.
func setupMiddleware(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(ApiTimeout))
	r.Use(middleware.AllowContentType("application/json", "charset=utf-8"))

	authHandler := authmiddleware.NewHandler()
	r.Use(authHandler.Handler)
}

// setupRoutes configures and registers all necessary routes with their respective handlers.
func setupRoutes(r *chi.Mux) {
	// Health check routes
	r.Get(ApiV1Prefix+"/health", healthhandler.HealthCheckHandler)

	// API key management routes
	r.Post(ApiV1Prefix+"/projects/{projectId}/api-keys", keyhandler.CreateApiKeyHandler)
	r.Get(ApiV1Prefix+"/projects/{projectId}/api-keys", keyhandler.ListApiKeysHandler)
	r.Get(ApiV1Prefix+"/projects/{projectId}/api-keys/{apiKey}", keyhandler.GetApiKeyHandler)
	r.Patch(ApiV1Prefix+"/projects/{projectId}/api-keys/{apiKey}", keyhandler.UpdateApiKeyHandler)
	r.Delete(ApiV1Prefix+"/projects/{projectId}/api-keys/{apiKey}", keyhandler.DeleteApiKeyHandler)

	// Prompt management routes
	r.Post(ApiV1Prefix+"/projects/{projectId}/prompts", prompthandler.CreatePromptHandler)
	r.Get(ApiV1Prefix+"/projects/{projectId}/prompts", prompthandler.ListPromptsHandler)

	r.Get(ApiV1Prefix+"/projects/{projectId}/prompts/{promptId}", prompthandler.GetPromptHandler)
	r.Patch(ApiV1Prefix+"/projects/{projectId}/prompts/{promptId}", prompthandler.UpdatePromptHandler)
	r.Delete(ApiV1Prefix+"/projects/{projectId}/prompts/{promptId}", prompthandler.DeletePromptHandler)

	r.Get(ApiV1Prefix+"/projects/{projectId}/prompts/{promptId}/versions", prompthandler.ListVersionsHandler)
	r.Put(ApiV1Prefix+"/projects/{projectId}/prompts/{promptId}/versions", prompthandler.UpdateLiveVersionHandler)

	r.Post(ApiV1Prefix+"/projects/{projectId}/prompts/{promptId}/branches", prompthandler.CreateBranchHandler)
	r.Get(ApiV1Prefix+"/projects/{projectId}/prompts/{promptId}/branches", prompthandler.ListBranchesHandler)
	r.Delete(ApiV1Prefix+"/projects/{projectId}/prompts/{promptId}/branches/{branchName}", prompthandler.DeleteBranchHandler)

	// Project management routes
	r.Get(ApiV1Prefix+"/projects", projecthandler.ListProjectsHandler)
	r.Post(ApiV1Prefix+"/projects", projecthandler.CreateProjectHandler)

	r.Get(ApiV1Prefix+"/projects/{projectId}", projecthandler.GetProjectHandler)
	r.Patch(ApiV1Prefix+"/projects/{projectId}", projecthandler.UpdateProjectHandler)
	r.Delete(ApiV1Prefix+"/projects/{projectId}", projecthandler.DeleteProjectHandler)

	// Team management routes
	r.Get(ApiV1Prefix+"/teams", teamhandler.ListTeamsHandler)
	r.Post(ApiV1Prefix+"/teams", teamhandler.CreateTeamHandler)

	r.Get(ApiV1Prefix+"/teams/{teamId}", teamhandler.GetTeamHandler)
	r.Patch(ApiV1Prefix+"/teams/{teamId}", teamhandler.UpdateTeamHandler)
	r.Delete(ApiV1Prefix+"/teams/{teamId}", teamhandler.DeleteTeamHandler)

	// User management routes
	r.Get(ApiV1Prefix+"/users", userhandler.ListUsersHandler)
	r.Post(ApiV1Prefix+"/users", userhandler.CreateUserHandler)

	r.Get(ApiV1Prefix+"/users/{userId}", userhandler.GetUserHandler)
	r.Patch(ApiV1Prefix+"/users/{userId}", userhandler.UpdateUserHandler)
	r.Delete(ApiV1Prefix+"/users/{userId}", userhandler.DeleteUserHandler)

	// Org management routes
	r.Post(ApiV1Prefix+"/orgs", orghandler.CreateOrgHandler)

	r.Get(ApiV1Prefix+"/orgs/{orgId}", orghandler.GetOrgHandler)
	r.Patch(ApiV1Prefix+"/orgs/{orgId}", orghandler.UpdateOrgHandler)
	r.Delete(ApiV1Prefix+"/orgs/{orgId}", orghandler.DeleteOrgHandler)
}
