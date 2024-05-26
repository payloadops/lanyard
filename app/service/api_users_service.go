package service

import (
	"context"
	"errors"
	"github.com/payloadops/plato/api/utils"
	"net/http"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
)

// UsersAPIService is a service that implements the logic for the UsersAPIServicer
// This service should implement the business logic for every endpoint for the UsersAPI API.
// Include any external packages or services that will be required by this service.
type UsersAPIService struct {
	client dal.UserManager
}

// NewUsersAPIService creates a default api service
func NewUsersAPIService() openapi.UsersAPIServicer {
	client, err := dal.NewUserDBClient()
	if err != nil {
		panic(err)
	}
	return &UsersAPIService{client: client}
}

// CreateUser - Create a new user
func (s *UsersAPIService) CreateUser(ctx context.Context, userInput openapi.UserInput) (openapi.ImplResponse, error) {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	user := dal.User{
		ID:    ksuid,
		Name:  userInput.Name,
		Email: userInput.Email,
	}

	err = s.client.CreateUser(ctx, user)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, user), nil
}

// DeleteUser - Delete a specific user
func (s *UsersAPIService) DeleteUser(ctx context.Context, userId string) (openapi.ImplResponse, error) {
	// Check if the user exists before trying to delete
	user, err := s.client.GetUser(ctx, userId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if user == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("user not found")
	}

	err = s.client.DeleteUser(ctx, userId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetUser - Get a specific user by ID
func (s *UsersAPIService) GetUser(ctx context.Context, userId string) (openapi.ImplResponse, error) {
	user, err := s.client.GetUser(ctx, userId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if user == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("user not found")
	}

	return openapi.Response(http.StatusOK, user), nil
}

// ListUsers - List all users
func (s *UsersAPIService) ListUsers(ctx context.Context) (openapi.ImplResponse, error) {
	users, err := s.client.ListUsers(ctx)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, users), nil
}

// UpdateUser - Update a specific user
func (s *UsersAPIService) UpdateUser(ctx context.Context, userId string, userInput openapi.UserInput) (openapi.ImplResponse, error) {
	// Check if the user exists before trying to update
	existingUser, err := s.client.GetUser(ctx, userId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if existingUser == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("user not found")
	}

	user := dal.User{
		ID:    userId,
		Name:  userInput.Name,
		Email: userInput.Email,
	}

	err = s.client.UpdateUser(ctx, user)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, user), nil
}
