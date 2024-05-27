package service

import (
	"context"
	"fmt"
	"github.com/payloadops/plato/api/utils"
	"net/http"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
)

// TeamsAPIService is a service that implements the logic for the TeamsAPIServicer
// This service should implement the business logic for every endpoint for the TeamsAPI API.
type TeamsAPIService struct {
	teamClient dal.TeamManager
	orgClient  dal.OrganizationManager
}

// NewTeamsAPIService creates a default app service
func NewTeamsAPIService() openapi.TeamsAPIServicer {
	teamClient, err := dal.NewTeamDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create TeamDBClient: %v", err))
	}
	orgClient, err := dal.NewOrgDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create OrgDBClient: %v", err))
	}
	return &TeamsAPIService{teamClient: teamClient, orgClient: orgClient}
}

// CreateTeam - Create a new team for an organization
func (s *TeamsAPIService) CreateTeam(ctx context.Context, orgId string, teamInput openapi.TeamInput) (openapi.ImplResponse, error) {
	// Check if the organization exists
	org, err := s.orgClient.GetOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if org == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("organization not found")
	}

	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	team := dal.Team{
		ID:          ksuid,
		OrgID:       orgId,
		Name:        teamInput.Name,
		Description: teamInput.Description,
	}

	err = s.teamClient.CreateTeam(ctx, team)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, team), nil
}

// DeleteTeam - Delete a specific team
func (s *TeamsAPIService) DeleteTeam(ctx context.Context, orgId string, teamId string) (openapi.ImplResponse, error) {
	// Check if the organization exists
	org, err := s.orgClient.GetOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if org == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("organization not found")
	}

	// Check if the team exists
	team, err := s.teamClient.GetTeam(ctx, teamId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if team == nil || team.OrgID != orgId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("team not found")
	}

	err = s.teamClient.DeleteTeam(ctx, teamId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetTeam - Get a specific team by ID
func (s *TeamsAPIService) GetTeam(ctx context.Context, orgId string, teamId string) (openapi.ImplResponse, error) {
	// Check if the organization exists
	org, err := s.orgClient.GetOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if org == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("organization not found")
	}

	team, err := s.teamClient.GetTeam(ctx, teamId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if team == nil || team.OrgID != orgId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("team not found")
	}

	return openapi.Response(http.StatusOK, team), nil
}

// ListTeams - List all teams in an organization
func (s *TeamsAPIService) ListTeams(ctx context.Context, orgId string) (openapi.ImplResponse, error) {
	// Check if the organization exists
	org, err := s.orgClient.GetOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if org == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("organization not found")
	}

	teams, err := s.teamClient.ListTeamsByOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, teams), nil
}

// UpdateTeam - Update a specific team
func (s *TeamsAPIService) UpdateTeam(ctx context.Context, orgId string, teamId string, teamInput openapi.TeamInput) (openapi.ImplResponse, error) {
	// Check if the organization exists
	org, err := s.orgClient.GetOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if org == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("organization not found")
	}

	// Check if the team exists
	team, err := s.teamClient.GetTeam(ctx, teamId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if team == nil || team.OrgID != orgId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("team not found")
	}

	// Update the team with the new values
	team.Name = teamInput.Name
	team.Description = teamInput.Description

	err = s.teamClient.UpdateTeam(ctx, *team)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, team), nil
}
