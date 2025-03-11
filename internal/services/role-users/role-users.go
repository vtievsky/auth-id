package roleusersvc

import (
	"context"
	"fmt"
	"time"

	"github.com/vtievsky/auth-id/internal/repositories/models"
	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
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

type Storage interface {
	GetRoleUsers(ctx context.Context, code string) ([]*models.RoleUser, error)
	AddRoleUser(ctx context.Context, roleUser models.RoleUserCreated) error
	UpdateRoleUser(ctx context.Context, roleUser models.RoleUserUpdated) error
	DeleteRoleUser(ctx context.Context, roleUser models.RoleUserDeleted) error
}

type UserSvc interface {
	GetUserByID(ctx context.Context, id uint64) (*usersvc.User, error)
	GetUserByLogin(ctx context.Context, login string) (*usersvc.User, error)
}

type RoleSvc interface {
	GetRoleByID(ctx context.Context, id uint64) (*rolesvc.Role, error)
	GetRoleByCode(ctx context.Context, code string) (*rolesvc.Role, error)
}

type RoleUserSvcOpts struct {
	Logger  *zap.Logger
	Storage Storage
	RoleSvc RoleSvc
	UserSvc UserSvc
}

type RoleUserSvc struct {
	logger  *zap.Logger
	storage Storage
	roleSvc RoleSvc
	userSvc UserSvc
}

func New(opts *RoleUserSvcOpts) *RoleUserSvc {
	return &RoleUserSvc{
		logger:  opts.Logger,
		storage: opts.Storage,
		roleSvc: opts.RoleSvc,
		userSvc: opts.UserSvc,
	}
}

func (s *RoleUserSvc) GetRoleUsers(ctx context.Context, code string) ([]*RoleUser, error) {
	const op = "RoleUserSvc.GetRoleUsers"

	ul, err := s.storage.GetRoleUsers(ctx, code)
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
				zap.Uint64("user_id", user.UserID),
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

func (s *RoleUserSvc) AddRoleUser(ctx context.Context, roleUser RoleUserCreated) error {
	const op = "RoleUserSvc.AddRoleUser"

	var (
		role *rolesvc.Role
		user *usersvc.User
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error

		role, err = s.roleSvc.GetRoleByCode(gCtx, roleUser.RoleCode)
		if err != nil {
			s.logger.Error("failed to parse role",
				zap.String("role_code", roleUser.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse role | %s:%w", op, err)
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

	if err := s.storage.AddRoleUser(ctx, models.RoleUserCreated{
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

func (s *RoleUserSvc) UpdateRoleUser(ctx context.Context, roleUser RoleUserUpdated) error {
	const op = "RoleUserSvc.UpdateRoleUser"

	var (
		role *rolesvc.Role
		user *usersvc.User
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error

		role, err = s.roleSvc.GetRoleByCode(gCtx, roleUser.RoleCode)
		if err != nil {
			s.logger.Error("failed to parse role",
				zap.String("role_code", roleUser.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse role | %s:%w", op, err)
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

	if err := s.storage.UpdateRoleUser(ctx, models.RoleUserUpdated{
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

func (s *RoleUserSvc) DeleteRoleUser(ctx context.Context, roleUser RoleUserDeleted) error {
	const op = "RoleUserSvc.DeleteRoleUser"

	var (
		role *rolesvc.Role
		user *usersvc.User
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error

		role, err = s.roleSvc.GetRoleByCode(gCtx, roleUser.RoleCode)
		if err != nil {
			s.logger.Error("failed to parse role",
				zap.String("role_code", roleUser.RoleCode),
				zap.Error(err),
			)

			return fmt.Errorf("failed to parse role | %s:%w", op, err)
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

	if err := s.storage.DeleteRoleUser(ctx, models.RoleUserDeleted{
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
