package redisusers

import (
	"context"
	"fmt"
	"time"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	redisclient "github.com/vtievsky/auth-id/internal/repositories/redis/client"
)

type UsersOpts struct {
	Client *redisclient.Client
}

type Users struct {
	client *redisclient.Client
}

func New(opts *UsersOpts) *Users {
	return &Users{
		client: opts.Client,
	}
}

func (s *Users) GetUser(ctx context.Context, login string) (*models.User, error) {
	const op = "DbUsers.GetUser"

	cmd := s.client.HGetAll(ctx, login)

	switch {
	case cmd.Err() != nil:
		return nil, fmt.Errorf("failed to get user | %s:%w", op, cmd.Err())
	case len(cmd.Val()) < 1:
		return nil, fmt.Errorf("failed to get user | %s:%w", op, dberrors.ErrUserNotFound)
	}

	var value redisclient.User

	err := cmd.Scan(&value)
	if err != nil {
		return nil, fmt.Errorf("failed to get user | %s:%w", op, dberrors.ErrUserScan)
	}

	return &models.User{
		ID:       value.ID,
		Login:    value.Login,
		FullName: value.FullName,
		Blocked:  value.Blocked,
	}, nil
}

func (s *Users) GetUsers(ctx context.Context) ([]*models.User, error) {
	users := make([]*models.User, 0)

	return users, nil
}

func (s *Users) CreateUser(ctx context.Context, user models.UserCreated) (*models.User, error) {
	const op = "DbUsers.CreateUser"

	if _, err := s.client.HMSet(ctx, user.Login, redisclient.UserCreated{
		ID:       int(time.Now().Unix()),
		Login:    user.Login,
		FullName: user.FullName,
		Blocked:  user.Blocked,
	}).Result(); err != nil {
		return nil, fmt.Errorf("failed to create user | %s:%w", op, err)
	}

	return s.GetUser(ctx, user.Login)
}

func (s *Users) UpdateUser(ctx context.Context, user models.UserUpdated) (*models.User, error) {
	const op = "DbUsers.UpdateUser"

	if _, err := s.client.HMSet(ctx, user.Login, redisclient.UserUpdated{
		Login:    user.Login,
		FullName: user.FullName,
		Blocked:  user.Blocked,
	}).Result(); err != nil {
		return nil, fmt.Errorf("failed to update user | %s:%w", op, err)
	}

	return s.GetUser(ctx, user.Login)
}

func (s *Users) DeleteUser(ctx context.Context, login string) error {
	const op = "DbUsers.DeleteUser"

	if _, err := s.client.Del(ctx, login).Result(); err != nil {
		return fmt.Errorf("failed to delete user | %s:%w", op, err)
	}

	return nil
}
