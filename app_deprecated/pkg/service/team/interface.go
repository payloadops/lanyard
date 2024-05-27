package team

import (
	"context"
	teamservicemodel "plato/app_deprecated/pkg/model/team"
)

type OrgService interface {
	CreateTeam(ctx context.Context, createTeamRequest teamservicemodel.CreateTeamRequest) (teamservicemodel.CreateTeamResponse, error)
	GetTeam(ctx context.Context, teamId string) (teamservicemodel.GetTeamResponse, error)
	ListTeams(ctx context.Context, orgId string) (teamservicemodel.ListTeamsResponse, error)
	UpdateTeam(ctx context.Context, projectId string, updateTeamRequest teamservicemodel.UpdateTeamRequest) (teamservicemodel.UpdateTeamResponse, error)
	DeleteTeam(ctx context.Context, projectId string) (teamservicemodel.DeleteTeamResponse, error)
}
