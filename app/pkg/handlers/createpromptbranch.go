package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func CreatePromptBranchHandler(w http.ResponseWriter, r *http.Request) {
	// validator := util.GetValidator()
	// if err := validator.Struct(createPromptBranchRequest); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	render.Status(r, http.StatusOK)
	// render.JSON(w, r, response)
}
