package rolesvc

import (
	"context"
	"fmt"
	"time"

	"github.com/vtievsky/auth-id/internal/repositories/models"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type RoleUser struct {
	Name    string
	Login   string
	DateIn  time.Time
	DateOut time.Time
}

type RoleUserCreated struct {
	Login    string
	RoleCode string
	DateIn   time.Time
	DateOut  time.Time
}

type RoleUserUpdated struct {
	Login    string
	RoleCode string
	DateIn   time.Time
	DateOut  time.Time
}

type RoleUserDeleted struct {
	Login    string
	RoleCode string
}

func (s *RoleSvc) GetRoleUsers(ctx context.Context, code string) ([]*RoleUser, error) {
	const op = "RoleSvc.GetRoleUsers"

	ul, err := s.roleUsers.GetRoleUsers(ctx, code)
	if err != nil {
		s.logger.Error("failed to get role users",
			zap.String("role_code", code),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get role users | %s:%w", op, err)
	}

	var p *usersvc.User

	resp := make([]*RoleUser, 0, len(ul))

	for _, user := range ul {
		p, err = s.userSvc.GetUserByID(ctx, user.UserID)
		if err != nil {
			s.logger.Error("failed to parse user",
				zap.String("role_code", code),
				zap.Int("user_id", user.UserID),
				zap.Error(err),
			)

			return nil, fmt.Errorf("failed to parse user | %s:%w", op, err)
		}

		resp = append(resp, &RoleUser{
			Login:   p.Login,
			Name:    p.Name,
			DateIn:  user.DateIn,
			DateOut: user.DateOut,
		})
	}

	return resp, nil
}

func (s *RoleSvc) AddRoleUser(ctx context.Context, roleUser RoleUserCreated) error {
	const op = "RoleSvc.AddRoleUser"

	var (
		role *models.Role
		user *usersvc.User
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error

		role, err = s.roles.GetRole(gCtx, roleUser.RoleCode)
		if err != nil {
			s.logger.Error("failed to get role",
				zap.String("role_code", roleUser.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to get role | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
		var err error

		user, err = s.userSvc.GetUserByLogin(gCtx, roleUser.Login)
		if err != nil {
			s.logger.Error("failed to parse user",
				zap.String("login", roleUser.Login),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse user | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if err := s.roleUsers.AddRoleUser(ctx, models.RoleUserCreated{
		RoleID:  role.ID,
		UserID:  user.ID,
		DateIn:  roleUser.DateIn,
		DateOut: roleUser.DateOut,
	}); err != nil {
		s.logger.Error("failed to add role to user",
			zap.String("role_code", roleUser.RoleCode),
			zap.String("login", roleUser.Login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to add role to user | %s:%w", op, err)
	}

	return nil
}

func (s *RoleSvc) UpdateRoleUser(ctx context.Context, roleUser RoleUserUpdated) error {
	const op = "RoleSvc.UpdateRoleUser"

	var (
		role *models.Role
		user *usersvc.User
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error

		role, err = s.roles.GetRole(gCtx, roleUser.RoleCode)
		if err != nil {
			s.logger.Error("failed to get role",
				zap.String("role_code", roleUser.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to get role | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
		var err error

		user, err = s.userSvc.GetUserByLogin(gCtx, roleUser.Login)
		if err != nil {
			s.logger.Error("failed to parse user",
				zap.String("login", roleUser.Login),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse user | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if err := s.roleUsers.UpdateRoleUser(ctx, models.RoleUserUpdated{
		RoleID:  role.ID,
		UserID:  user.ID,
		DateIn:  roleUser.DateIn,
		DateOut: roleUser.DateOut,
	}); err != nil {
		s.logger.Error("failed to update role to user",
			zap.String("role_code", roleUser.RoleCode),
			zap.String("login", roleUser.Login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to update role to user | %s:%w", op, err)
	}

	return nil
}

func (s *RoleSvc) DeleteRoleUser(ctx context.Context, roleUser RoleUserDeleted) error {
	const op = "RoleSvc.DeleteRoleUser"

	var (
		role *models.Role
		user *usersvc.User
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error

		role, err = s.roles.GetRole(gCtx, roleUser.RoleCode)
		if err != nil {
			s.logger.Error("failed to get role",
				zap.String("role_code", roleUser.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to get role | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
		var err error

		user, err = s.userSvc.GetUserByLogin(gCtx, roleUser.Login)
		if err != nil {
			s.logger.Error("failed to parse user",
				zap.String("login", roleUser.Login),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse user | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if err := s.roleUsers.DeleteRoleUser(ctx, models.RoleUserDeleted{
		RoleID: role.ID,
		UserID: user.ID,
	}); err != nil {
		s.logger.Error("failed to delete role to user",
			zap.String("role_code", roleUser.RoleCode),
			zap.String("login", roleUser.Login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to delete role to user | %s:%w", op, err)
	}

	return nil
}
