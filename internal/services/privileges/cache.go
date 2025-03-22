package privilegesvc

import (
	"context"
	"fmt"
	"math"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	"go.uber.org/zap"
)

func (s *PrivilegeSvc) GetPrivilegeByID(ctx context.Context, id uint64) (*Privilege, error) {
	const op = "PrivilegeSvc.GetPrivilegeByID"

	val, err := s.cacheByID.Get(ctx, id, s.syncPrivilegesByID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w | %v", op, dberrors.ErrPrivilegeNotFound, err)
	}

	return &Privilege{
		ID:          val.ID,
		Code:        val.Code,
		Name:        val.Name,
		Description: val.Description,
	}, nil
}

func (s *PrivilegeSvc) GetPrivilegeByCode(ctx context.Context, code string) (*Privilege, error) {
	const op = "PrivilegeSvc.GetPrivilegeByCode"

	val, err := s.cacheByCode.Get(ctx, code, s.syncPrivilegesByCode)
	if err != nil {
		return nil, fmt.Errorf("%s:%w | %v", op, dberrors.ErrPrivilegeNotFound, err)
	}

	return &Privilege{
		ID:          val.ID,
		Code:        val.Code,
		Name:        val.Name,
		Description: val.Description,
	}, nil
}

func (s *PrivilegeSvc) syncPrivilegesByID(ctx context.Context) (map[uint64]*models.Privilege, error) {
	const op = "PrivilegeSvc.syncPrivilegesByID"

	privileges, err := s.storage.GetPrivileges(ctx, math.MaxUint32, 0)
	if err != nil {
		s.logger.Error("failed to sync privileges",
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to sync privileges | %s:%w", op, err)
	}

	resp := make(map[uint64]*models.Privilege, len(privileges))

	for _, privilege := range privileges {
		resp[privilege.ID] = privilege
	}

	s.logger.Debug("privileges has been synchronized",
		zap.Int("num", len(privileges)),
	)

	return resp, nil
}

func (s *PrivilegeSvc) syncPrivilegesByCode(ctx context.Context) (map[string]*models.Privilege, error) {
	const op = "PrivilegeSvc.syncPrivilegesByCode"

	privileges, err := s.storage.GetPrivileges(ctx, math.MaxUint32, 0)
	if err != nil {
		s.logger.Error("failed to sync privileges",
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to sync privileges | %s:%w", op, err)
	}

	resp := make(map[string]*models.Privilege, len(privileges))

	for _, privilege := range privileges {
		resp[privilege.Code] = privilege
	}

	s.logger.Debug("privileges has been synchronized",
		zap.Int("num", len(privileges)),
	)

	return resp, nil
}
