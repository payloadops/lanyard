package model

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	StatusCode int    `json:"-"`
	StatusText string `json:"status"`
	Message    string `json:"error,omitempty"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func ErrorResponseRenderer(statusCode int, message string) render.Renderer {
	return &ErrorResponse{
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
		Message:    message,
	}
}
