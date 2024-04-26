package user

import "context"

type OrgService interface {
	CreateUser(ctx context.Context, createUserRequest userservicemodel.CreateUserRequest) (userservicemodel.CreateUserResponse, error)
	GetUser(ctx context.Context, userId string) (userservicemodel.GetUserResponse, error)
	ListUsersByOrg(ctx context.Context, orgId string) (userservicemodel.ListUserResponse, error)
	ListUsersByTeam(ctx context.Context, teamId string) (userservicemodel.ListUserResponse, error)
	UpdateUser(ctx context.Context, projectId string, updateUserRequest userservicemodel.UpdateUserRequest) (userservicemodel.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, projectId string) (userservicemodel.DeleteUserResponse, error)
}
