package model

import (
	"bytes"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorResponse_Render(t *testing.T) {
	// Create an instance of ErrorResponse
	errResponse := &ErrorResponse{
		StatusCode: http.StatusBadRequest,
		StatusText: http.StatusText(http.StatusBadRequest),
		Message:    "Bad request",
	}

	// Create a mock HTTP response writer
	w := httptest.NewRecorder()

	// Create a mock HTTP request
	r := httptest.NewRequest("GET", "/foo", bytes.NewReader([]byte("{}")))

	// Call the Render method
	err := errResponse.Render(w, r)

	// Check if Render method returns no error
	assert.NoError(t, err)

	// Set the status code and render the error into the body
	render.Render(w, r, errResponse)

	// Check if the status code is correctly set
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check if the response body contains the expected JSON
	expectedJSON := `{"status":"Bad Request","error":"Bad request"}`
	assert.JSONEq(t, expectedJSON, w.Body.String())
}

func TestErrorResponseRenderer(t *testing.T) {
	// Create an ErrorResponse using ErrorResponseRenderer
	statusCode := http.StatusNotFound
	message := "Resource not found"
	errResponse := ErrorResponseRenderer(statusCode, message)

	// Check if the status code is correctly set
	assert.Equal(t, statusCode, errResponse.(*ErrorResponse).StatusCode)

	// Check if the status text is correctly set
	assert.Equal(t, http.StatusText(statusCode), errResponse.(*ErrorResponse).StatusText)

	// Check if the message is correctly set
	assert.Equal(t, message, errResponse.(*ErrorResponse).Message)
}
