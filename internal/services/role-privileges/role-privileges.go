package roleprivilegesvc

import (
	"context"
	"fmt"

	"github.com/vtievsky/auth-id/internal/repositories/models"
	privilegesvc "github.com/vtievsky/auth-id/internal/services/privileges"
	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type RolePrivilege struct {
	Code        string
	Name        string
	Description string
	Allowed     bool
}

type RolePrivilegeCreated struct {
	RoleCode      string
	PrivilegeCode string
	Allowed       bool
}

type RolePrivilegeUpdated struct {
	RoleCode      string
	PrivilegeCode string
	Allowed       bool
}

type RolePrivilegeDeleted struct {
	RoleCode      string
	PrivilegeCode string
}

type Storage interface {
	GetRolePrivileges(ctx context.Context, code string, pageSize, offset uint32) ([]*models.RolePrivilege, error)
	AddRolePrivilege(ctx context.Context, rolePrivilege models.RolePrivilegeCreated) error
	UpdateRolePrivilege(ctx context.Context, rolePrivilege models.RolePrivilegeUpdated) error
	DeleteRolePrivilege(ctx context.Context, rolePrivilege models.RolePrivilegeDeleted) error
}

type RoleSvc interface {
	GetRoleByID(ctx context.Context, id uint64) (*rolesvc.Role, error)
	GetRoleByCode(ctx context.Context, code string) (*rolesvc.Role, error)
}

type PrivilegeSvc interface {
	GetPrivilegeByID(ctx context.Context, id uint64) (*privilegesvc.Privilege, error)
	GetPrivilegeByCode(ctx context.Context, code string) (*privilegesvc.Privilege, error)
}

type RolePrivilegeSvcOpts struct {
	Logger       *zap.Logger
	Storage      Storage
	RoleSvc      RoleSvc
	PrivilegeSvc PrivilegeSvc
}

type RolePrivilegeSvc struct {
	logger       *zap.Logger
	storage      Storage
	roleSvc      RoleSvc
	privilegeSvc PrivilegeSvc
}

func New(opts *RolePrivilegeSvcOpts) *RolePrivilegeSvc {
	return &RolePrivilegeSvc{
		logger:       opts.Logger,
		storage:      opts.Storage,
		roleSvc:      opts.RoleSvc,
		privilegeSvc: opts.PrivilegeSvc,
	}
}

func (s *RolePrivilegeSvc) GetRolePrivileges(ctx context.Context, code string, pageSize, offset uint32) ([]*RolePrivilege, error) {
	const op = "RolePrivilegeSvc.GetRolePrivileges"

	ul, err := s.storage.GetRolePrivileges(ctx, code, pageSize, offset)
	if err != nil {
		s.logger.Error("failed to get role privileges",
			zap.String("role_code", code),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get role privileges | %s:%w", op, err)
	}

	var p *privilegesvc.Privilege

	resp := make([]*RolePrivilege, 0, len(ul))

	for _, privilege := range ul {
		p, err = s.privilegeSvc.GetPrivilegeByID(ctx, privilege.PrivilegeID)
		if err != nil {
			s.logger.Error("failed to parse privilege",
				zap.String("role_code", code),
				zap.Uint64("privilege_id", privilege.PrivilegeID),
				zap.Error(err),
			)

			return nil, fmt.Errorf("failed to parse privilege | %s:%w", op, err)
		}

		resp = append(resp, &RolePrivilege{
			Code:        p.Code,
			Name:        p.Name,
			Description: p.Description,
			Allowed:     privilege.Allowed,
		})
	}

	return resp, nil
}

func (s *RolePrivilegeSvc) AddRolePrivilege(ctx context.Context, rolePrivilege RolePrivilegeCreated) error {
	const op = "RolePrivilegeSvc.AddRolePrivilege"

	var (
		role      *rolesvc.Role
		privilege *privilegesvc.Privilege
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error

		role, err = s.roleSvc.GetRoleByCode(gCtx, rolePrivilege.RoleCode)
		if err != nil {
			s.logger.Error("failed to parse role",
				zap.String("role_code", rolePrivilege.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse role | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
		var err error

		privilege, err = s.privilegeSvc.GetPrivilegeByCode(gCtx, rolePrivilege.PrivilegeCode)
		if err != nil {
			s.logger.Error("failed to parse privilege",
				zap.String("privilege_code", rolePrivilege.PrivilegeCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse privilege | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if err := s.storage.AddRolePrivilege(ctx, models.RolePrivilegeCreated{
		RoleID:      role.ID,
		PrivilegeID: privilege.ID,
		Allowed:     rolePrivilege.Allowed,
	}); err != nil {
		s.logger.Error("failed to add role to privilege",
			zap.String("role_code", rolePrivilege.RoleCode),
			zap.String("privilege_code", rolePrivilege.PrivilegeCode),
			zap.Error(err),
		)

		return fmt.Errorf("failed to add role to privilege | %s:%w", op, err)
	}

	return nil
}

func (s *RolePrivilegeSvc) UpdateRolePrivilege(ctx context.Context, rolePrivilege RolePrivilegeUpdated) error {
	const op = "RolePrivilegeSvc.UpdateRolePrivilege"

	var (
		role      *rolesvc.Role
		privilege *privilegesvc.Privilege
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error

		role, err = s.roleSvc.GetRoleByCode(gCtx, rolePrivilege.RoleCode)
		if err != nil {
			s.logger.Error("failed to parse role",
				zap.String("role_code", rolePrivilege.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse role | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
		var err error

		privilege, err = s.privilegeSvc.GetPrivilegeByCode(gCtx, rolePrivilege.PrivilegeCode)
		if err != nil {
			s.logger.Error("failed to parse privilege",
				zap.String("privilege_code", rolePrivilege.PrivilegeCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse privilege | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if err := s.storage.UpdateRolePrivilege(ctx, models.RolePrivilegeUpdated{
		RoleID:      role.ID,
		PrivilegeID: privilege.ID,
		Allowed:     rolePrivilege.Allowed,
	}); err != nil {
		s.logger.Error("failed to update role to privilege",
			zap.String("role_code", rolePrivilege.RoleCode),
			zap.String("privilege_code", rolePrivilege.PrivilegeCode),
			zap.Error(err),
		)

		return fmt.Errorf("failed to update role to privilege | %s:%w", op, err)
	}

	return nil
}

func (s *RolePrivilegeSvc) DeleteRolePrivilege(ctx context.Context, rolePrivilege RolePrivilegeDeleted) error {
	const op = "RolePrivilegeSvc.DeleteRolePrivilege"

	var (
		role      *rolesvc.Role
		privilege *privilegesvc.Privilege
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error

		role, err = s.roleSvc.GetRoleByCode(gCtx, rolePrivilege.RoleCode)
		if err != nil {
			s.logger.Error("failed to parse role",
				zap.String("role_code", rolePrivilege.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse role | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
		var err error

		privilege, err = s.privilegeSvc.GetPrivilegeByCode(gCtx, rolePrivilege.PrivilegeCode)
		if err != nil {
			s.logger.Error("failed to parse privilege",
				zap.String("privilege_code", rolePrivilege.PrivilegeCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse privilege | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if err := s.storage.DeleteRolePrivilege(ctx, models.RolePrivilegeDeleted{
		RoleID:      role.ID,
		PrivilegeID: privilege.ID,
	}); err != nil {
		s.logger.Error("failed to delete role to privilege",
			zap.String("role_code", rolePrivilege.RoleCode),
			zap.String("privilege_code", rolePrivilege.PrivilegeCode),
			zap.Error(err),
		)

		return fmt.Errorf("failed to delete role to privilege | %s:%w", op, err)
	}

	return nil
}
