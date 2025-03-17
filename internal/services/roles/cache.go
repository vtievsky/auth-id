package rolesvc

import (
	"context"
	"fmt"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	"go.uber.org/zap"
)

func (s *RoleSvc) GetRoleByID(ctx context.Context, id uint64) (*Role, error) {
	const op = "RoleSvc.GetRoleByID"

	val, err := s.cacheByID.Get(ctx, id, s.syncRolesByID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w | %v", op, dberrors.ErrRoleNotFound, err)
	}

	return &Role{
		ID:          val.ID,
		Code:        val.Code,
		Name:        val.Name,
		Description: val.Description,
		Blocked:     val.Blocked,
	}, nil
}

func (s *RoleSvc) GetRoleByCode(ctx context.Context, code string) (*Role, error) {
	const op = "RoleSvc.GetRoleByCode"

	val, err := s.cacheByCode.Get(ctx, code, s.syncRolesByCode)
	if err != nil {
		return nil, fmt.Errorf("%s:%w | %v", op, dberrors.ErrRoleNotFound, err)
	}

	return &Role{
		ID:          val.ID,
		Code:        val.Code,
		Name:        val.Name,
		Description: val.Description,
		Blocked:     val.Blocked,
	}, nil
}

func (s *RoleSvc) syncRolesByID(ctx context.Context) (map[uint64]*models.Role, error) {
	const op = "RoleSvc.syncRolesByID"

	roles, err := s.storage.GetRoles(ctx)
	if err != nil {
		s.logger.Error("failed to sync roles",
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to sync roles | %s:%w", op, err)
	}

	resp := make(map[uint64]*models.Role, len(roles))

	for _, role := range roles {
		resp[role.ID] = role
	}

	s.logger.Debug("roles has been synchronized",
		zap.Int("num", len(roles)),
	)

	return resp, nil
}

func (s *RoleSvc) syncRolesByCode(ctx context.Context) (map[string]*models.Role, error) {
	const op = "RoleSvc.syncRolesByCode"

	roles, err := s.storage.GetRoles(ctx)
	if err != nil {
		s.logger.Error("failed to sync roles",
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to sync roles | %s:%w", op, err)
	}

	resp := make(map[string]*models.Role, len(roles))

	for _, role := range roles {
		resp[role.Code] = role
	}

	s.logger.Debug("roles has been synchronized",
		zap.Int("num", len(roles)),
	)

	return resp, nil
}
