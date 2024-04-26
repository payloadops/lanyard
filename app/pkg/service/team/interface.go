package team

import "context"

type OrgService interface {
	CreateTeam(ctx context.Context, createTeamRequest model.CreateTeamRequest) (model.CreateTeamResponse, error)
	GetTeam(ctx context.Context, teamId string) (model.GetTeamResponse, error)
	ListTeams(ctx context.Context, orgId string) (model.ListTeamResponse, error)
	UpdateTeam(ctx context.Context, projectId string, updateTeamRequest model.UpdateTeamRequest) (model.UpdateTeamResponse, error)
	DeleteTeam(ctx context.Context, projectId string) (model.DeleteTeamResponse, error)
}
