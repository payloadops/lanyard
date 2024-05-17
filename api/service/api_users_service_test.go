package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserManager is a mock implementation of the UserManager interface
type MockUserManager struct {
	mock.Mock
}

func (m *MockUserManager) CreateUser(ctx context.Context, user dal.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserManager) GetUser(ctx context.Context, id string) (*dal.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dal.User), args.Error(1)
}

func (m *MockUserManager) UpdateUser(ctx context.Context, user dal.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserManager) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserManager) ListUsers(ctx context.Context) ([]dal.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dal.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockUserManager := new(MockUserManager)
	service := UsersAPIService{client: mockUserManager}

	userInput := openapi.UserInput{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mockUserManager.On("CreateUser", mock.Anything, dal.User{
		Name:  userInput.Name,
		Email: userInput.Email,
	}).Return(nil)

	resp, err := service.CreateUser(context.Background(), userInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	mockUserManager.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockUserManager := new(MockUserManager)
	service := UsersAPIService{client: mockUserManager}

	userID := "1"

	// Test case where user exists
	mockUserManager.On("GetUser", mock.Anything, userID).Return(&dal.User{}, nil)
	mockUserManager.On("DeleteUser", mock.Anything, userID).Return(nil)

	resp, err := service.DeleteUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockUserManager.AssertExpectations(t)

	// Test case where user does not exist
	mockUserManager.On("GetUser", mock.Anything, userID).Return((*dal.User)(nil), nil)

	resp, err = service.DeleteUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockUserManager.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	mockUserManager := new(MockUserManager)
	service := UsersAPIService{client: mockUserManager}

	userID := "1"
	user := &dal.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Test case where user exists
	mockUserManager.On("GetUser", mock.Anything, userID).Return(user, nil)

	resp, err := service.GetUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockUserManager.AssertExpectations(t)

	// Test case where user does not exist
	mockUserManager.On("GetUser", mock.Anything, userID).Return((*dal.User)(nil), nil)

	resp, err = service.GetUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockUserManager.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
	mockUserManager := new(MockUserManager)
	service := UsersAPIService{client: mockUserManager}

	users := []dal.User{
		{ID: "1", Name: "John Doe", Email: "john@example.com"},
		{ID: "2", Name: "Jane Doe", Email: "jane@example.com"},
	}

	mockUserManager.On("ListUsers", mock.Anything).Return(users, nil)

	resp, err := service.ListUsers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockUserManager.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockUserManager := new(MockUserManager)
	service := UsersAPIService{client: mockUserManager}

	userID := "1"
	userInput := openapi.UserInput{
		Name:  "John Doe Updated",
		Email: "john.updated@example.com",
	}
	user := dal.User{
		ID:    userID,
		Name:  userInput.Name,
		Email: userInput.Email,
	}

	// Test case where user exists
	mockUserManager.On("GetUser", mock.Anything, userID).Return(&dal.User{}, nil)
	mockUserManager.On("UpdateUser", mock.Anything, user).Return(nil)

	resp, err := service.UpdateUser(context.Background(), userID, userInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockUserManager.AssertExpectations(t)

	// Test case where user does not exist
	mockUserManager.On("GetUser", mock.Anything, userID).Return((*dal.User)(nil), nil)

	resp, err = service.UpdateUser(context.Background(), userID, userInput)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockUserManager.AssertExpectations(t)
}
