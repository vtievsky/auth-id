package userrolesvc

import (
	"context"
	"fmt"
	"sync"
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

type Users interface {
	GetUser(ctx context.Context, login string) (*models.User, error)
}

type UserRoles interface {
	GetUserRoles(ctx context.Context, login string) ([]*models.UserRole, error)
}

type RoleSvc interface {
	GetRoleByID(ctx context.Context, id int) (*rolesvc.Role, error)
	GetRoleByCode(ctx context.Context, code string) (*rolesvc.Role, error)
}

type UserRoleSvcOpts struct {
	Logger    *zap.Logger
	Users     Users
	UserRoles UserRoles
	RoleSvc   RoleSvc
}

type UserRoleSvc struct {
	logger       *zap.Logger
	users        Users
	userRoles    UserRoles
	roleSvc      RoleSvc
	lastTime     time.Time
	cacheByID    map[int]*models.User
	cacheByLogin map[string]*models.User
	mu           sync.RWMutex
}

func New(opts *UserRoleSvcOpts) *UserRoleSvc {
	return &UserRoleSvc{
		logger:       opts.Logger,
		users:        opts.Users,
		userRoles:    opts.UserRoles,
		roleSvc:      opts.RoleSvc,
		lastTime:     time.Time{},
		cacheByID:    make(map[int]*models.User),
		cacheByLogin: make(map[string]*models.User),
		mu:           sync.RWMutex{},
	}
}

func (s *UserRoleSvc) GetUserRoles(ctx context.Context, login string) ([]*UserRole, error) {
	const op = "UserRoleSvc.GetUserRoles"

	ul, err := s.userRoles.GetUserRoles(ctx, login)
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
				zap.Int("role_id", role.RoleID),
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
