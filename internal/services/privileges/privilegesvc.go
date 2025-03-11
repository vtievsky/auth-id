package privilegesvc

import (
	"context"

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
