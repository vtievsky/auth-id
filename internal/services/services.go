package services

import (
	"context"

	privilegesvc "github.com/vtievsky/auth-id/internal/services/privileges"
	roleprivilegesvc "github.com/vtievsky/auth-id/internal/services/role-privileges"
	roleusersvc "github.com/vtievsky/auth-id/internal/services/role-users"
	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
	sessionsvc "github.com/vtievsky/auth-id/internal/services/sessions"
	userprivilegesvc "github.com/vtievsky/auth-id/internal/services/user-privileges"
	userrolesvc "github.com/vtievsky/auth-id/internal/services/user-roles"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
)

type SvcLayer struct {
	UserSvc          UserService
	UserRoleSvc      UserRoleService
	UserPrivilegeSvc UserPrivilegeService
	RoleSvc          RoleService
	RoleUserSvc      RoleUserService
	RolePrivilegeSvc RolePrivilegeService
	PrivilegeSvc     PrivilegeService
	SessionSvc       SessionService
}

type UserService interface {
	GetUser(ctx context.Context, login string) (*usersvc.User, error)
	GetUsers(ctx context.Context, pageSize, offset uint32) ([]*usersvc.User, error)
	CreateUser(ctx context.Context, user usersvc.UserCreated) (*usersvc.User, error)
	UpdateUser(ctx context.Context, user usersvc.UserUpdated) (*usersvc.User, error)
	DeleteUser(ctx context.Context, login string) error

	ChangePass(ctx context.Context, login, current, changed string) error
	ResetPass(ctx context.Context, login, changed string) error
}

type UserRoleService interface {
	GetUserRoles(ctx context.Context, login string, pageSize, offset uint32) ([]*userrolesvc.UserRole, error)
}

type UserPrivilegeService interface {
	GetUserPrivileges(ctx context.Context, login string, pageSize, offset uint32) ([]*userprivilegesvc.UserPrivilege, error)
}

type RoleService interface {
	GetRole(ctx context.Context, code string) (*rolesvc.Role, error)
	GetRoles(ctx context.Context, pageSize, offset uint32) ([]*rolesvc.Role, error)
	CreateRole(ctx context.Context, user rolesvc.RoleCreated) (*rolesvc.Role, error)
	UpdateRole(ctx context.Context, user rolesvc.RoleUpdated) (*rolesvc.Role, error)
	DeleteRole(ctx context.Context, code string) error
}

type RoleUserService interface {
	GetRoleUsers(ctx context.Context, code string, pageSize, offset uint32) ([]*roleusersvc.RoleUser, error)
	AddRoleUser(ctx context.Context, roleUser roleusersvc.RoleUserCreated) error
	UpdateRoleUser(ctx context.Context, roleUser roleusersvc.RoleUserUpdated) error
	DeleteRoleUser(ctx context.Context, roleUser roleusersvc.RoleUserDeleted) error
}

type RolePrivilegeService interface {
	GetRolePrivileges(ctx context.Context, code string, pageSize, offset uint32) ([]*roleprivilegesvc.RolePrivilege, error)
	AddRolePrivilege(ctx context.Context, rolePrivilege roleprivilegesvc.RolePrivilegeCreated) error
	UpdateRolePrivilege(ctx context.Context, rolePrivilege roleprivilegesvc.RolePrivilegeUpdated) error
	DeleteRolePrivilege(ctx context.Context, rolePrivilege roleprivilegesvc.RolePrivilegeDeleted) error
}

type SessionService interface {
	Get(ctx context.Context, sessionID string) (*sessionsvc.SessionCart, error)
	Login(ctx context.Context, login, password string) (*sessionsvc.Tokens, error)
	GetUserSessions(ctx context.Context, login string, pageSize, offset uint32) ([]*sessionsvc.Session, error)
	Delete(ctx context.Context, login, sessionID string) error
	Search(ctx context.Context, sessionID, privilege string) error
}

type PrivilegeService interface {
	GetPrivileges(ctx context.Context, pageSize, offset uint32) ([]*privilegesvc.Privilege, error)
}
