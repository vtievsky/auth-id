package privilegesvc

import (
	"context"
	"fmt"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"go.uber.org/zap"
)

func (s *PrivilegeSvc) GetPrivilegeByID(ctx context.Context, id uint64) (*Privilege, error) {
	const op = "PrivilegeSvc.GetPrivilegeByID"

	val, err := s.cacheByID.Get(ctx, id, s.syncPrivileges)
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

	val, err := s.cacheByCode.Get(ctx, code, s.syncPrivileges)
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

func (s *PrivilegeSvc) syncPrivileges(ctx context.Context) error {
	const op = "PrivilegeSvc.syncPrivileges"

	resp, err := s.storage.GetPrivileges(ctx)
	if err != nil {
		s.logger.Error("failed to sync privileges",
			zap.Error(err),
		)

		return fmt.Errorf("failed to sync privileges | %s:%w", op, err)
	}

	for _, privilege := range resp {
		s.cacheByID.Add(privilege.ID, privilege)
		s.cacheByCode.Add(privilege.Code, privilege)
	}

	return nil
}
