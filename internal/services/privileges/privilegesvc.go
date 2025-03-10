package privilegesvc

import (
	"context"
	"fmt"
	"sync"
	"time"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	"go.uber.org/zap"
)

const (
	cacheTTL = time.Second * 60
)

type Privilege struct {
	ID          int    // Ключевое поле для БД
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
	lastTime    time.Time
	cacheByID   map[int]*models.Privilege
	cacheByCode map[string]*models.Privilege
	mu          sync.RWMutex
}

func New(opts *PrivilegeSvcOpts) *PrivilegeSvc {
	return &PrivilegeSvc{
		logger:      opts.Logger,
		storage:     opts.Storage,
		lastTime:    time.Time{},
		cacheByID:   make(map[int]*models.Privilege),
		cacheByCode: make(map[string]*models.Privilege),
		mu:          sync.RWMutex{},
	}
}

func (s *PrivilegeSvc) GetPrivilegeByID(ctx context.Context, id int) (*Privilege, error) {
	const op = "PrivilegeSvc.GetPrivilegeByID"

	s.mu.RLock()

	// Кеш невалидный
	if cacheTTL < time.Since(s.lastTime) {
		s.mu.RLocker().Unlock()

		if err := s.syncPrivileges(ctx); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		s.mu.RLock()
	}

	defer s.mu.RLocker().Unlock()

	if val, ok := s.cacheByID[id]; ok {
		return &Privilege{
			ID:          val.ID,
			Code:        val.Code,
			Name:        val.Name,
			Description: val.Description,
		}, nil
	}

	return nil, dberrors.ErrPrivilegeNotFound
}

func (s *PrivilegeSvc) GetPrivilegeByCode(ctx context.Context, code string) (*Privilege, error) {
	const op = "PrivilegeSvc.GetPrivilegeByCode"

	s.mu.RLock()

	// Кеш невалидный
	if cacheTTL < time.Since(s.lastTime) {
		s.mu.RLocker().Unlock()

		if err := s.syncPrivileges(ctx); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		s.mu.RLock()
	}

	defer s.mu.RLocker().Unlock()

	if val, ok := s.cacheByCode[code]; ok {
		return &Privilege{
			ID:          val.ID,
			Code:        val.Code,
			Name:        val.Name,
			Description: val.Description,
		}, nil
	}

	return nil, dberrors.ErrPrivilegeNotFound
}

func (s *PrivilegeSvc) syncPrivileges(ctx context.Context) error {
	const op = "PrivilegeSvc.syncPrivileges"

	s.mu.Lock()
	defer s.mu.Unlock()

	resp, err := s.storage.GetPrivileges(ctx)
	if err != nil {
		s.logger.Error("failed to sync privileges",
			zap.Error(err),
		)

		return fmt.Errorf("failed to sync privileges | %s:%w", op, err)
	}

	for _, privilege := range resp {
		s.cacheByID[privilege.ID] = privilege
		s.cacheByCode[privilege.Code] = privilege
	}

	// Зафиксируем время синхронизации справочника
	s.lastTime = time.Now()

	return nil
}
