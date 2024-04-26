package team

import "context"

type OrgService interface {
	CreateTeam(ctx context.Context, createTeamRequest teamservicemodel.CreateTeamRequest) (teamservicemodel.CreateTeamResponse, error)
	GetTeam(ctx context.Context, teamId string) (teamservicemodel.GetTeamResponse, error)
	ListTeams(ctx context.Context, orgId string) (teamservicemodel.ListTeamResponse, error)
	UpdateTeam(ctx context.Context, projectId string, updateTeamRequest teamservicemodel.UpdateTeamRequest) (teamservicemodel.UpdateTeamResponse, error)
	DeleteTeam(ctx context.Context, projectId string) (teamservicemodel.DeleteTeamResponse, error)
}
