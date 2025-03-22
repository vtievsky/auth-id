package privilegesvc

import (
	"context"
	"fmt"

	"github.com/vtievsky/auth-id/internal/repositories/models"
	"github.com/vtievsky/auth-id/pkg/cache"
	"go.uber.org/zap"
)

type Privilege struct {
	ID          uint64 // Ключевое поле для БД
	Code        string // Ключевое поле для интерфейса
	Name        string
	Description string
}

type Storage interface {
	GetPrivileges(ctx context.Context) ([]*models.Privilege, error)
}

type PrivilegeSvcOpts struct {
	Logger  *zap.Logger
	Storage Storage
}

type PrivilegeSvc struct {
	logger      *zap.Logger
	storage     Storage
	cacheByID   cache.Cache[uint64, *models.Privilege]
	cacheByCode cache.Cache[string, *models.Privilege]
}

func New(opts *PrivilegeSvcOpts) *PrivilegeSvc {
	return &PrivilegeSvc{
		logger:      opts.Logger,
		storage:     opts.Storage,
		cacheByID:   cache.New[uint64, *models.Privilege](),
		cacheByCode: cache.New[string, *models.Privilege](),
	}
}

func (s *PrivilegeSvc) GetPrivileges(ctx context.Context) ([]*Privilege, error) {
	const op = "PrivilegeSvc.GetPrivileges"

	privileges, err := s.storage.GetPrivileges(ctx)
	if err != nil {
		s.logger.Error("failed to get privileges",
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get privileges | %s:%w", op, err)
	}

	resp := make([]*Privilege, 0, len(privileges))

	for _, privilege := range privileges {
		resp = append(resp, &Privilege{
			ID:          privilege.ID,
			Code:        privilege.Code,
			Name:        privilege.Name,
			Description: privilege.Description,
		})
	}

	return resp, nil
}
