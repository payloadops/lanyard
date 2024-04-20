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
