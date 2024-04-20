package router

import (
	"net/http"
	"plato/app/pkg/handlers"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Init() {
	r := chi.NewRouter()

	// Initialize middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Setup routes
	setupRoutes(r)

	// Start server
	http.ListenAndServe(":8080", r)
}

func setupRoutes(r *chi.Mux) {
	r.Get("/health", handlers.HealthCheckHandler)

	r.Post("/prompts", handlers.CreatePromptHandler)
	r.Get("/prompts", handlers.ListPromptsHandler)

	r.Get("/prompts/{promptId}", handlers.GetPromptHandler)
	r.Patch("/prompts/{promptId}", handlers.UpdatePromptHandler)
	r.Delete("/prompts/{promptId}", handlers.DeletePromptHandler)

	r.Get("/prompts/{promptId}/versions", handlers.ListPromptVersionsHandler)
	r.Put("/prompts/{promptId}/versions/{versionId}", handlers.UpdatePromptVersionHandler)

	r.Post("/prompts/{promptId}/branches", handlers.CreatePromptBranchHandler)
	r.Delete("/prompts/{promptId}/branches/{branchId}", handlers.DeletePromptBranchHandler)
}
