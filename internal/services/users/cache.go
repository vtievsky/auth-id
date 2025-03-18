package usersvc

import (
	"context"
	"fmt"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	"go.uber.org/zap"
)

func (s *UserSvc) GetUserByID(ctx context.Context, id uint64) (*User, error) {
	const op = "UserSvc.GetUserByID"

	val, err := s.cacheByID.Get(ctx, id, s.syncUsersByID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w | %v", op, dberrors.ErrUserNotFound, err)
	}

	return &User{
		ID:       val.ID,
		Name:     val.Name,
		Login:    val.Login,
		Password: val.Password,
		Blocked:  val.Blocked,
	}, nil
}

func (s *UserSvc) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	const op = "UserSvc.GetUserByLogin"

	val, err := s.cacheByLogin.Get(ctx, login, s.syncUsersByLogin)
	if err != nil {
		s.logger.Error("failed to get user",
			zap.Error(err),
		)

		return nil, fmt.Errorf("%s:%w | %v", op, dberrors.ErrUserNotFound, err)
	}

	return &User{
		ID:       val.ID,
		Name:     val.Name,
		Login:    val.Login,
		Password: val.Password,
		Blocked:  val.Blocked,
	}, nil
}

func (s *UserSvc) syncUsersByID(ctx context.Context) (map[uint64]*models.User, error) {
	const op = "UserSvc.syncUsersByID"

	users, err := s.storage.GetUsers(ctx)
	if err != nil {
		s.logger.Error("failed to sync users",
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to sync users | %s:%w", op, err)
	}

	resp := make(map[uint64]*models.User, len(users))

	for _, user := range users {
		resp[user.ID] = user
	}

	s.logger.Debug("users has been synchronized",
		zap.Int("num", len(users)),
	)

	return resp, nil
}

func (s *UserSvc) syncUsersByLogin(ctx context.Context) (map[string]*models.User, error) {
	const op = "UserSvc.syncUsersByLogin"

	users, err := s.storage.GetUsers(ctx)
	if err != nil {
		s.logger.Error("failed to sync users",
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to sync users | %s:%w", op, err)
	}

	resp := make(map[string]*models.User, len(users))

	for _, user := range users {
		resp[user.Login] = user
	}

	s.logger.Debug("users has been synchronized",
		zap.Int("num", len(users)),
	)

	return resp, nil
}
