package handlers

import (
	"errors"
	"net/http"
)

func validateListPromptsRequest() error {
	if false {
		return errors.New("name is required")
	}
	return nil
}

func ListPromptsHandler(w http.ResponseWriter, r *http.Request) {

}
