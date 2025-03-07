package services

import (
	"context"

	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
)

type SvcLayer struct {
	UserSvc UserService
	RoleSvc RoleService
}

type UserService interface {
	GetUser(ctx context.Context, login string) (*usersvc.User, error)
	GetUsers(ctx context.Context) ([]*usersvc.User, error)
	CreateUser(ctx context.Context, user usersvc.UserCreated) (*usersvc.User, error)
	UpdateUser(ctx context.Context, user usersvc.UserUpdated) (*usersvc.User, error)
	DeleteUser(ctx context.Context, login string) error
}

type RoleService interface {
	GetRole(ctx context.Context, id int) (*rolesvc.Role, error)
	GetRoles(ctx context.Context) ([]*rolesvc.Role, error)
	CreateRole(ctx context.Context, user rolesvc.RoleCreated) (*rolesvc.Role, error)
	UpdateRole(ctx context.Context, user rolesvc.RoleUpdated) (*rolesvc.Role, error)
	DeleteRole(ctx context.Context, id int) error
}
