package rolesvc

import (
	"context"
	"fmt"

	"github.com/vtievsky/auth-id/internal/repositories/models"
	privilegesvc "github.com/vtievsky/auth-id/internal/services/privileges"
	"go.uber.org/zap"
)

type Role struct {
	ID          int
	Code        string
	Name        string
	Description string
	Blocked     bool
}

type RoleCreated struct {
	Name        string
	Description string
	Blocked     bool
}

type RoleUpdated struct {
	Code        string
	Name        string
	Description string
	Blocked     bool
}

type Roles interface {
	GetRole(ctx context.Context, code string) (*models.Role, error)
	GetRoles(ctx context.Context) ([]*models.Role, error)
	CreateRole(ctx context.Context, user models.RoleCreated) (*models.Role, error)
	UpdateRole(ctx context.Context, user models.RoleUpdated) (*models.Role, error)
	DeleteRole(ctx context.Context, code string) error
}

type RolePrivileges interface {
	GetRolePrivileges(ctx context.Context, code string) ([]*models.RolePrivilege, error)
	AddRolePrivilege(ctx context.Context, rolePrivilege models.RolePrivilegeCreated) error
}

type PrivilegeSvc interface {
	GetPrivilegeByID(ctx context.Context, id int) (*privilegesvc.Privilege, error)
	GetPrivilegeByCode(ctx context.Context, code string) (*privilegesvc.Privilege, error)
}

type RoleSvcOpts struct {
	Logger         *zap.Logger
	Roles          Roles
	RolePrivileges RolePrivileges
	PrivilegeSvc   PrivilegeSvc
}

type RoleSvc struct {
	logger         *zap.Logger
	roles          Roles
	rolePrivileges RolePrivileges
	privilegeSvc   PrivilegeSvc
}

func New(opts *RoleSvcOpts) *RoleSvc {
	return &RoleSvc{
		logger:         opts.Logger,
		roles:          opts.Roles,
		rolePrivileges: opts.RolePrivileges,
		privilegeSvc:   opts.PrivilegeSvc,
	}
}

func (s *RoleSvc) GetRole(ctx context.Context, code string) (*Role, error) {
	const op = "RoleSvc.GetUser"

	resp, err := s.roles.GetRole(ctx, code)
	if err != nil {
		s.logger.Error("failed to get role",
			zap.String("role_code", code),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get role | %s:%w", op, err)
	}

	return &Role{
		ID:          resp.ID,
		Code:        resp.Code,
		Name:        resp.Name,
		Description: resp.Description,
		Blocked:     false,
	}, nil
}

func (s *RoleSvc) GetRoles(ctx context.Context) ([]*Role, error) {
	const op = "RoleSvc.GetRoles"

	ul, err := s.roles.GetRoles(ctx)
	if err != nil {
		s.logger.Error("failed to get roles",
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get roles | %s:%w", op, err)
	}

	roles := make([]*Role, 0, len(ul))

	for _, role := range ul {
		roles = append(roles, &Role{
			ID:          role.ID,
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			Blocked:     role.Blocked,
		})
	}

	return roles, nil
}

func (s *RoleSvc) CreateRole(ctx context.Context, role RoleCreated) (*Role, error) {
	const op = "RoleSvc.CreateRole"

	u, err := s.roles.CreateRole(ctx, models.RoleCreated{
		Name:        role.Name,
		Description: role.Description,
		Blocked:     role.Blocked,
	})
	if err != nil {
		s.logger.Error("failed to create role",
			zap.String("role_name", role.Name),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to create role | %s:%w", op, err)
	}

	return &Role{
		ID:          u.ID,
		Code:        u.Code,
		Name:        u.Name,
		Description: u.Description,
		Blocked:     u.Blocked,
	}, nil
}

func (s *RoleSvc) UpdateRole(ctx context.Context, role RoleUpdated) (*Role, error) {
	const op = "RoleSvc.UpdateRole"

	u, err := s.roles.UpdateRole(ctx, models.RoleUpdated{
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		Blocked:     role.Blocked,
	})
	if err != nil {
		s.logger.Error("failed to update role",
			zap.String("role_code", role.Code),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to update role | %s:%w", op, err)
	}

	return &Role{
		ID:          u.ID,
		Code:        u.Code,
		Name:        u.Name,
		Description: u.Description,
		Blocked:     u.Blocked,
	}, nil
}

func (s *RoleSvc) DeleteRole(ctx context.Context, code string) error {
	const op = "RoleSvc.RoleUser"

	if err := s.roles.DeleteRole(ctx, code); err != nil {
		s.logger.Error("failed to delete role",
			zap.String("role_code", code),
			zap.Error(err),
		)

		return fmt.Errorf("failed to delete role | %s:%w", op, err)
	}

	return nil
}
