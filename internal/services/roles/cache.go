package rolesvc

import (
	"context"
	"fmt"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"go.uber.org/zap"
)

func (s *RoleSvc) GetRoleByID(ctx context.Context, id uint64) (*Role, error) {
	const op = "RoleSvc.GetRoleByID"

	val, err := s.cacheByID.Get(ctx, id, s.syncRoles)
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

	val, err := s.cacheByCode.Get(ctx, code, s.syncRoles)
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

func (s *RoleSvc) syncRoles(ctx context.Context) error {
	const op = "RoleSvc.syncRoles"

	resp, err := s.storage.GetRoles(ctx)
	if err != nil {
		s.logger.Error("failed to sync roles",
			zap.Error(err),
		)

		return fmt.Errorf("failed to sync roles | %s:%w", op, err)
	}

	for _, role := range resp {
		s.cacheByID.Add(role.ID, role)
		s.cacheByCode.Add(role.Code, role)
	}

	s.logger.Debug("roles has been synchronized",
		zap.Int("num", len(resp)),
	)

	return nil
}
