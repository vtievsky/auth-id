package rolesvc

import (
	"context"
	"fmt"

	"github.com/vtievsky/auth-id/internal/repositories/models"
	privilegesvc "github.com/vtievsky/auth-id/internal/services/privileges"
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

func (s *RoleSvc) GetRolePrivileges(ctx context.Context, code string) ([]*RolePrivilege, error) {
	const op = "RoleSvc.GetRolePrivileges"

	ul, err := s.rolePrivileges.GetRolePrivileges(ctx, code)
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
				zap.Int("privilege_id", privilege.PrivilegeID),
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

func (s *RoleSvc) AddRolePrivilege(ctx context.Context, rolePrivilege RolePrivilegeCreated) error {
	const op = "RoleSvc.AddRolePrivilege"

	var (
		err       error
		role      *models.Role
		privilege *privilegesvc.Privilege
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		role, err = s.roles.GetRole(gCtx, rolePrivilege.RoleCode)
		if err != nil {
			s.logger.Error("failed to get role",
				zap.String("role_code", rolePrivilege.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to get role | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
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

	if err = g.Wait(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if err := s.rolePrivileges.AddRolePrivilege(ctx, models.RolePrivilegeCreated{
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

func (s *RoleSvc) UpdateRolePrivilege(ctx context.Context, rolePrivilege RolePrivilegeUpdated) error {
	const op = "RoleSvc.UpdateRolePrivilege"

	var (
		err       error
		role      *models.Role
		privilege *privilegesvc.Privilege
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		role, err = s.roles.GetRole(gCtx, rolePrivilege.RoleCode)
		if err != nil {
			s.logger.Error("failed to get role",
				zap.String("role_code", rolePrivilege.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to get role | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
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

	if err = g.Wait(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if err := s.rolePrivileges.UpdateRolePrivilege(ctx, models.RolePrivilegeUpdated{
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

func (s *RoleSvc) DeleteRolePrivilege(ctx context.Context, rolePrivilege RolePrivilegeDeleted) error {
	const op = "RoleSvc.DeleteRolePrivilege"

	var (
		err       error
		role      *models.Role
		privilege *privilegesvc.Privilege
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		role, err = s.roles.GetRole(gCtx, rolePrivilege.RoleCode)
		if err != nil {
			s.logger.Error("failed to get role",
				zap.String("role_code", rolePrivilege.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to get role | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
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

	if err = g.Wait(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if err := s.rolePrivileges.DeleteRolePrivilege(ctx, models.RolePrivilegeDeleted{
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
