package rolesvc

import (
	"context"
	"fmt"

	"github.com/vtievsky/auth-id/internal/repositories/models"
	"github.com/vtievsky/auth-id/pkg/cache"
	"go.uber.org/zap"
)

type Role struct {
	ID          uint64
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

type Storage interface {
	GetRole(ctx context.Context, code string) (*models.Role, error)
	GetRoles(ctx context.Context, pageSize, offset uint32) ([]*models.Role, error)
	CreateRole(ctx context.Context, user models.RoleCreated) (*models.Role, error)
	UpdateRole(ctx context.Context, user models.RoleUpdated) (*models.Role, error)
	DeleteRole(ctx context.Context, code string) error
}

type RoleSvcOpts struct {
	Logger  *zap.Logger
	Storage Storage
}

type RoleSvc struct {
	logger      *zap.Logger
	storage     Storage
	cacheByID   cache.Cache[uint64, *models.Role]
	cacheByCode cache.Cache[string, *models.Role]
}

func New(opts *RoleSvcOpts) *RoleSvc {
	return &RoleSvc{
		logger:      opts.Logger,
		storage:     opts.Storage,
		cacheByID:   cache.New[uint64, *models.Role](),
		cacheByCode: cache.New[string, *models.Role](),
	}
}

func (s *RoleSvc) GetRole(ctx context.Context, code string) (*Role, error) {
	return s.GetRoleByCode(ctx, code)
}

func (s *RoleSvc) GetRoles(ctx context.Context, pageSize, offset uint32) ([]*Role, error) {
	const op = "RoleSvc.GetRoles"

	ul, err := s.storage.GetRoles(ctx, pageSize, offset)
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

	u, err := s.storage.CreateRole(ctx, models.RoleCreated{
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

	s.cacheByID.Add(u.ID, u)
	s.cacheByCode.Add(u.Code, u)

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

	u, err := s.storage.UpdateRole(ctx, models.RoleUpdated{
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
	const op = "RoleSvc.DeleteRole"

	u, err := s.GetRoleByCode(ctx, code)
	if err != nil {
		s.logger.Error("failed to get role",
			zap.String("role_code", code),
			zap.Error(err),
		)

		return fmt.Errorf("failed to get role | %s:%w", op, err)
	}

	if err := s.storage.DeleteRole(ctx, code); err != nil {
		s.logger.Error("failed to delete role",
			zap.String("role_code", code),
			zap.Error(err),
		)

		return fmt.Errorf("failed to delete role | %s:%w", op, err)
	}

	// Удалим данные роли из кеша
	s.cacheByID.Del(u.ID)
	s.cacheByCode.Del(u.Code)

	return nil
}
