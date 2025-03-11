package usersvc

import (
	"context"
	"fmt"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"go.uber.org/zap"
)

func (s *UserSvc) GetUserByID(ctx context.Context, id uint64) (*User, error) {
	const op = "UserSvc.GetUserByID"

	val, err := s.cacheByID.Get(ctx, id, s.syncUsers)
	if err != nil {
		return nil, fmt.Errorf("%s:%w | %v", op, dberrors.ErrUserNotFound, err)
	}

	return &User{
		ID:      val.ID,
		Login:   val.Login,
		Name:    val.Name,
		Blocked: val.Blocked,
	}, nil
}

func (s *UserSvc) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	const op = "UserSvc.GetUserByLogin"

	val, err := s.cacheByLogin.Get(ctx, login, s.syncUsers)
	if err != nil {
		return nil, fmt.Errorf("%s:%w | %v", op, dberrors.ErrUserNotFound, err)
	}

	return &User{
		ID:      val.ID,
		Login:   val.Login,
		Name:    val.Name,
		Blocked: val.Blocked,
	}, nil
}

func (s *UserSvc) syncUsers(ctx context.Context) error {
	const op = "UserSvc.syncUsers"

	resp, err := s.storage.GetUsers(ctx)
	if err != nil {
		s.logger.Error("failed to sync users",
			zap.Error(err),
		)

		return fmt.Errorf("failed to sync users | %s:%w", op, err)
	}

	for _, user := range resp {
		s.cacheByID.Add(user.ID, user)
		s.cacheByLogin.Add(user.Login, user)
	}

	return nil
}
