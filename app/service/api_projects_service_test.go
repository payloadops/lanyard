package service

import (
	"context"
	"go.uber.org/mock/gomock"
	"testing"
	"time"

	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/dal/mocks"
	"github.com/payloadops/plato/app/openapi"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func TestProjectsAPIService_CreateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockProjectManager(ctrl)
	service := NewProjectsAPIService(mockClient)

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectInput := openapi.ProjectInput{
		Name:        "Project1",
		Description: "Description1",
	}

	mockClient.EXPECT().CreateProject(ctx, "org1", gomock.Any()).Return(nil)

	response, err := service.CreateProject(ctx, projectInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.NotNil(t, response.Body)
	project, ok := response.Body.(openapi.Project)
	assert.True(t, ok)
	assert.Equal(t, projectInput.Name, project.Name)
	assert.Equal(t, projectInput.Description, project.Description)
}

func TestProjectsAPIService_DeleteProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockProjectManager(ctrl)
	service := NewProjectsAPIService(mockClient)

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"

	mockClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockClient.EXPECT().DeleteProject(ctx, "org1", projectID).Return(nil)

	response, err := service.DeleteProject(ctx, projectID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestProjectsAPIService_GetProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockProjectManager(ctrl)
	service := NewProjectsAPIService(mockClient)

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"

	mockClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{
		ProjectID:   projectID,
		Name:        "Project1",
		Description: "Description1",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}, nil)

	response, err := service.GetProject(ctx, projectID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	project, ok := response.Body.(openapi.Project)
	assert.True(t, ok)
	assert.Equal(t, projectID, project.Id)
}

func TestProjectsAPIService_ListProjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockProjectManager(ctrl)
	service := NewProjectsAPIService(mockClient)

	ctx := context.WithValue(context.Background(), "orgID", "org1")

	projects := []dal.Project{
		{ProjectID: "proj1", Name: "Project1", Description: "Description1", CreatedAt: time.Now().UTC().Format(time.RFC3339), UpdatedAt: time.Now().UTC().Format(time.RFC3339)},
		{ProjectID: "proj2", Name: "Project2", Description: "Description2", CreatedAt: time.Now().UTC().Format(time.RFC3339), UpdatedAt: time.Now().UTC().Format(time.RFC3339)},
	}

	mockClient.EXPECT().ListProjectsByOrganization(ctx, "org1").Return(projects, nil)

	response, err := service.ListProjects(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	listedProjects, ok := response.Body.([]openapi.Project)
	assert.True(t, ok)
	assert.Equal(t, 2, len(listedProjects))
	assert.Equal(t, "proj1", listedProjects[0].Id)
	assert.Equal(t, "proj2", listedProjects[1].Id)
}

func TestProjectsAPIService_UpdateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockProjectManager(ctrl)
	service := NewProjectsAPIService(mockClient)

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	projectInput := openapi.ProjectInput{
		Name:        "UpdatedProject",
		Description: "UpdatedDescription",
	}

	mockClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{
		ProjectID:   projectID,
		Name:        "Project1",
		Description: "Description1",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}, nil)
	mockClient.EXPECT().UpdateProject(ctx, "org1", gomock.Any()).Return(nil)

	response, err := service.UpdateProject(ctx, projectID, projectInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	updatedProject, ok := response.Body.(openapi.Project)
	assert.True(t, ok)
	assert.Equal(t, projectInput.Name, updatedProject.Name)
	assert.Equal(t, projectInput.Description, updatedProject.Description)
}
