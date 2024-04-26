package user

import "context"

type OrgService interface {
	CreateUser(ctx context.Context, createUserRequest model.CreateUserRequest) (model.CreateUserResponse, error)
	GetUser(ctx context.Context, userId string) (model.GetUserResponse, error)
	ListUsersByOrg(ctx context.Context, orgId string) (model.ListUserResponse, error)
	ListUsersByTeam(ctx context.Context, teamId string) (model.ListUserResponse, error)
	UpdateUser(ctx context.Context, projectId string) error
	DeleteUser(ctx context.Context, projectId string) error
}
