package rolesvc

import (
	"context"
	"fmt"
	"time"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"go.uber.org/zap"
)

const (
	cacheTTL = time.Second * 60
)

func (s *RoleSvc) GetRoleByID(ctx context.Context, id uint64) (*Role, error) {
	const op = "RoleSvc.GetRoleByID"

	s.mu.RLock()

	// Кеш невалидный
	if cacheTTL < time.Since(s.lastTime) {
		s.mu.RLocker().Unlock()

		if err := s.syncRoles(ctx); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		s.mu.RLock()
	}

	defer s.mu.RLocker().Unlock()

	if val, ok := s.cacheByID[id]; ok {
		return &Role{
			ID:          val.ID,
			Code:        val.Code,
			Name:        val.Name,
			Description: val.Description,
			Blocked:     val.Blocked,
		}, nil
	}

	return nil, dberrors.ErrRoleNotFound
}

func (s *RoleSvc) GetRoleByCode(ctx context.Context, code string) (*Role, error) {
	const op = "RoleSvc.GetRoleByCode"

	s.mu.RLock()

	// Кеш невалидный
	if cacheTTL < time.Since(s.lastTime) {
		s.mu.RLocker().Unlock()

		if err := s.syncRoles(ctx); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		s.mu.RLock()
	}

	defer s.mu.RLocker().Unlock()

	if val, ok := s.cacheByCode[code]; ok {
		return &Role{
			ID:          val.ID,
			Code:        val.Code,
			Name:        val.Name,
			Description: val.Description,
			Blocked:     val.Blocked,
		}, nil
	}

	return nil, dberrors.ErrUserNotFound
}

func (s *RoleSvc) syncRoles(ctx context.Context) error {
	const op = "RoleSvc.syncRoles"

	s.mu.Lock()
	defer s.mu.Unlock()

	resp, err := s.roles.GetRoles(ctx)
	if err != nil {
		s.logger.Error("failed to sync roles",
			zap.Error(err),
		)

		return fmt.Errorf("failed to sync roles | %s:%w", op, err)
	}

	for _, role := range resp {
		s.cacheByID[role.ID] = role
		s.cacheByCode[role.Code] = role
	}

	// Зафиксируем время синхронизации справочника
	s.lastTime = time.Now()

	return nil
}
