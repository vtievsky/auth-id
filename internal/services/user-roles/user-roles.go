package userrolesvc

import (
	"context"
	"fmt"
	"time"

	"github.com/vtievsky/auth-id/internal/repositories/models"
	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
	"go.uber.org/zap"
)

type UserRole struct {
	Code        string
	Name        string
	Description string
	DateIn      time.Time
	DateOut     time.Time
}

type Storage interface {
	GetUserRoles(ctx context.Context, login string, pageSize, offset uint32) ([]*models.UserRole, error)
}

type RoleSvc interface {
	GetRoleByID(ctx context.Context, id uint64) (*rolesvc.Role, error)
	GetRoleByCode(ctx context.Context, code string) (*rolesvc.Role, error)
}

type UserRoleSvcOpts struct {
	Logger  *zap.Logger
	Storage Storage
	RoleSvc RoleSvc
}

type UserRoleSvc struct {
	logger  *zap.Logger
	storage Storage
	roleSvc RoleSvc
}

func New(opts *UserRoleSvcOpts) *UserRoleSvc {
	return &UserRoleSvc{
		logger:  opts.Logger,
		storage: opts.Storage,
		roleSvc: opts.RoleSvc,
	}
}

func (s *UserRoleSvc) GetUserRoles(ctx context.Context, login string, pageSize, offset uint32) ([]*UserRole, error) {
	const op = "UserRoleSvc.GetUserRoles"

	ul, err := s.storage.GetUserRoles(ctx, login, pageSize, offset)
	if err != nil {
		s.logger.Error("failed to get user roles",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get user roles | %s:%w", op, err)
	}

	var p *rolesvc.Role

	resp := make([]*UserRole, 0, len(ul))

	for _, role := range ul {
		p, err = s.roleSvc.GetRoleByID(ctx, role.RoleID)
		if err != nil {
			s.logger.Error("failed to parse role",
				zap.String("login", login),
				zap.Uint64("role_id", role.RoleID),
				zap.Error(err),
			)

			return nil, fmt.Errorf("failed to parse role | %s:%w", op, err)
		}

		resp = append(resp, &UserRole{
			Code:        p.Code,
			Name:        p.Name,
			Description: p.Description,
			DateIn:      role.DateIn,
			DateOut:     role.DateOut,
		})
	}

	return resp, nil
}
