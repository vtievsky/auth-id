package redisusers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UsersOpts struct {
	URL string
}

type Users struct {
	client redis.UniversalClient
}

func New(opts *UsersOpts) *Users {
	client := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs:      []string{opts.URL},
			ClientName: "auth-id",
			DB:         0,
			PoolSize:   20,
		},
	)

	return &Users{
		client: client,
	}
}

func (s *Users) GetUser(ctx context.Context, login string) (*User, error) {
	const op = "DbUsers.GetUser"

	var value User

	if err := s.client.HGetAll(ctx, login).Scan(&value); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("failed to get user | %s:%w", op, ErrUserNotFound)
		}

		return nil, fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	return &User{
		ID:       value.ID,
		Login:    value.Login,
		FullName: value.FullName,
		Blocked:  value.Blocked,
	}, nil
}

func (s *Users) GetUsers(ctx context.Context) ([]*User, error) {
	users := make([]*User, 0)

	return users, nil
}

func (s *Users) CreateUser(ctx context.Context, user UserCreated) (*User, error) {
	const op = "DbUsers.CreateUser"

	if _, err := s.client.HMSet(ctx, user.Login, UserCreated{
		ID:       int(time.Now().Unix()),
		Login:    user.Login,
		FullName: user.FullName,
		Blocked:  user.Blocked,
	}).Result(); err != nil {
		return nil, fmt.Errorf("failed to create user | %s:%w", op, err)
	}

	return s.GetUser(ctx, user.Login)
}

func (s *Users) UpdateUser(ctx context.Context, user UserUpdated) (*User, error) {
	const op = "DbUsers.UpdateUser"

	if _, err := s.client.HMSet(ctx, user.Login, UserUpdated{
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
