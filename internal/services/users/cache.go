package usersvc

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

func (s *UserSvc) GetUserByID(ctx context.Context, id uint64) (*User, error) {
	const op = "UserSvc.GetUserByID"

	s.mu.RLock()

	// Кеш невалидный
	if cacheTTL < time.Since(s.lastTime) {
		s.mu.RLocker().Unlock()

		if err := s.syncUsers(ctx); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		s.mu.RLock()
	}

	defer s.mu.RLocker().Unlock()

	if val, ok := s.cacheByID[id]; ok {
		return &User{
			ID:      val.ID,
			Login:   val.Login,
			Name:    val.Name,
			Blocked: val.Blocked,
		}, nil
	}

	return nil, dberrors.ErrUserNotFound
}

func (s *UserSvc) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	const op = "UserSvc.GetUserByLogin"

	s.mu.RLock()

	// Кеш невалидный
	if cacheTTL < time.Since(s.lastTime) {
		s.mu.RLocker().Unlock()

		if err := s.syncUsers(ctx); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		s.mu.RLock()
	}

	defer s.mu.RLocker().Unlock()

	if val, ok := s.cacheByLogin[login]; ok {
		return &User{
			ID:      val.ID,
			Login:   val.Login,
			Name:    val.Name,
			Blocked: val.Blocked,
		}, nil
	}

	return nil, dberrors.ErrUserNotFound
}

func (s *UserSvc) syncUsers(ctx context.Context) error {
	const op = "UserSvc.syncUsers"

	s.mu.Lock()
	defer s.mu.Unlock()

	resp, err := s.users.GetUsers(ctx)
	if err != nil {
		s.logger.Error("failed to sync users",
			zap.Error(err),
		)

		return fmt.Errorf("failed to sync users | %s:%w", op, err)
	}

	for _, user := range resp {
		s.cacheByID[user.ID] = user
		s.cacheByLogin[user.Login] = user
	}

	// Зафиксируем время синхронизации справочника
	s.lastTime = time.Now()

	return nil
}
