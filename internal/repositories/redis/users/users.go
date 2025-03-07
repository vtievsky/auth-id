package redisusers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	redisclient "github.com/vtievsky/auth-id/internal/repositories/redis/client"
)

const (
	space = "usr"
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

	cmd := s.client.HGetAll(ctx, s.loginToKey(login))

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
		ID:      value.ID,
		Login:   value.Login,
		Name:    value.Name,
		Blocked: value.Blocked,
	}, nil
}

func (s *Users) GetUsers(ctx context.Context) ([]*models.User, error) {
	const op = "DbUsers.GetUsers"

	ul, err := s.client.Keys(ctx, s.space()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get users | %s:%w", op, dberrors.ErrUserScan)
	}

	users := make([]*models.User, 0)

	for _, key := range ul {
		u, err := s.GetUser(ctx, s.keyToLogin(key))
		if err != nil {
			continue
		}

		users = append(users, &models.User{
			ID:      u.ID,
			Login:   u.Login,
			Name:    u.Name,
			Blocked: u.Blocked,
		})
	}

	return users, nil
}

func (s *Users) CreateUser(ctx context.Context, user models.UserCreated) (*models.User, error) {
	const op = "DbUsers.CreateUser"

	if _, err := s.GetUser(ctx, user.Login); err != nil {
		if !errors.Is(err, dberrors.ErrUserNotFound) {
			return nil, fmt.Errorf("failed to create user | %s:%w", op, err)
		}

		userID := int(time.Now().Unix())

		if _, err := s.client.HMSet(ctx, s.loginToKey(user.Login), redisclient.UserCreated{
			ID:      userID,
			Login:   user.Login,
			Name:    user.Name,
			Blocked: user.Blocked,
		}).Result(); err != nil {
			return nil, fmt.Errorf("failed to create user | %s:%w", op, err)
		}

		return s.GetUser(ctx, user.Login)
	}

	return nil, fmt.Errorf("failed to create user | %s:%w", op, dberrors.ErrUserAlreadyExists)
}

func (s *Users) UpdateUser(ctx context.Context, user models.UserUpdated) (*models.User, error) {
	const op = "DbUsers.UpdateUser"

	if _, err := s.GetUser(ctx, user.Login); err != nil {
		return nil, fmt.Errorf("failed to update user | %s:%w", op, err)
	}

	if _, err := s.client.HMSet(ctx, s.loginToKey(user.Login), redisclient.UserUpdated{
		Login:   user.Login,
		Name:    user.Name,
		Blocked: user.Blocked,
	}).Result(); err != nil {
		return nil, fmt.Errorf("failed to update user | %s:%w", op, err)
	}

	return s.GetUser(ctx, user.Login)
}

func (s *Users) DeleteUser(ctx context.Context, login string) error {
	const op = "DbUsers.DeleteUser"

	if _, err := s.GetUser(ctx, login); err != nil {
		return fmt.Errorf("failed to delete user | %s:%w", op, err)
	}

	if _, err := s.client.Del(ctx, s.loginToKey(login)).Result(); err != nil {
		return fmt.Errorf("failed to delete user | %s:%w", op, err)
	}

	return nil
}

func (s *Users) space() string {
	return fmt.Sprintf("%s:*", space)
}

func (s *Users) loginToKey(login string) string {
	return fmt.Sprintf("%s:%s", space, login)
}

func (s *Users) keyToLogin(key string) string {
	p := fmt.Sprintf("%s:", space)

	return strings.ReplaceAll(key, p, "")
}
