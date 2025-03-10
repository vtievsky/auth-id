package services

import (
	"context"

	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
)

type SvcLayer struct {
	UserSvc          UserService
	RoleSvc          RoleService
	RolePrivilegeSvc RolePrivilegeService
	RoleUserSvc      RoleUserService
}

type UserService interface {
	GetUser(ctx context.Context, login string) (*usersvc.User, error)
	GetUsers(ctx context.Context) ([]*usersvc.User, error)
	CreateUser(ctx context.Context, user usersvc.UserCreated) (*usersvc.User, error)
	UpdateUser(ctx context.Context, user usersvc.UserUpdated) (*usersvc.User, error)
	DeleteUser(ctx context.Context, login string) error
}

type RoleService interface {
	GetRole(ctx context.Context, code string) (*rolesvc.Role, error)
	GetRoles(ctx context.Context) ([]*rolesvc.Role, error)
	CreateRole(ctx context.Context, user rolesvc.RoleCreated) (*rolesvc.Role, error)
	UpdateRole(ctx context.Context, user rolesvc.RoleUpdated) (*rolesvc.Role, error)
	DeleteRole(ctx context.Context, code string) error
}

type RolePrivilegeService interface {
	GetRolePrivileges(ctx context.Context, code string) ([]*rolesvc.RolePrivilege, error)
	AddRolePrivilege(ctx context.Context, rolePrivilege rolesvc.RolePrivilegeCreated) error
	UpdateRolePrivilege(ctx context.Context, rolePrivilege rolesvc.RolePrivilegeUpdated) error
	DeleteRolePrivilege(ctx context.Context, rolePrivilege rolesvc.RolePrivilegeDeleted) error
}

type RoleUserService interface {
	GetRoleUsers(ctx context.Context, code string) ([]*rolesvc.RoleUser, error)
	AddRoleUser(ctx context.Context, roleUser rolesvc.RoleUserCreated) error
	UpdateRoleUser(ctx context.Context, roleUser rolesvc.RoleUserUpdated) error
	DeleteRoleUser(ctx context.Context, roleUser rolesvc.RoleUserDeleted) error
}
