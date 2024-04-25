package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func DeletePromptBranchHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	// render.JSON(w, r, response)
}
