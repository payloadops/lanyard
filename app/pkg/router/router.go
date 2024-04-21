package router

import (
	"net/http"
	"plato/app/pkg/handlers"
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
}

func setupRoutes(r *chi.Mux) {
	r.Get(API_V1_PREFIX+"/health", handlers.HealthCheckHandler)

	r.Post(API_V1_PREFIX+"/prompts", handlers.CreatePromptHandler)
	r.Get(API_V1_PREFIX+"/prompts", handlers.ListPromptsHandler)

	r.Get(API_V1_PREFIX+"/prompts/{promptId}", handlers.GetPromptHandler)
	r.Patch(API_V1_PREFIX+"/prompts/{promptId}", handlers.UpdatePromptHandler)
	r.Delete(API_V1_PREFIX+"/prompts/{promptId}", handlers.DeletePromptHandler)

	r.Get(API_V1_PREFIX+"/prompts/{promptId}/versions", handlers.ListPromptVersionsHandler)
	r.Put(API_V1_PREFIX+"/prompts/{promptId}/versions/{versionId}", handlers.UpdatePromptVersionHandler)

	r.Post(API_V1_PREFIX+"/prompts/{promptId}/branches", handlers.CreatePromptBranchHandler)
	r.Delete(API_V1_PREFIX+"/prompts/{promptId}/branches/{branchId}", handlers.DeletePromptBranchHandler)
}
