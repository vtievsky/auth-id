package tarantoolusers

import (
	"context"
	"errors"
	"fmt"

	"github.com/tarantool/go-tarantool"
	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	tarantoolclient "github.com/vtievsky/auth-id/internal/repositories/tarantool/client"
)

const (
	space = "user"
	limit = 25
)

type UsersOpts struct {
	Client *tarantoolclient.Client
}

type Users struct {
	c *tarantoolclient.Client
}

func New(opts *UsersOpts) *Users {
	return &Users{
		c: opts.Client,
	}
}

func (s *Users) GetUser(ctx context.Context, login string) (*models.User, error) {
	const op = "DbUsers.GetUser"

	resp, err := s.c.Connection.Select(space, "secondary", 0, 1, tarantool.IterEq, tarantoolclient.Tuple{login})
	if err != nil {
		return nil, fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	if len(resp.Tuples()) < 1 {
		return nil, fmt.Errorf("failed to get user | %s:%w", op, dberrors.ErrUserNotFound)
	}

	value := s.tupleToUser(resp.Tuples()[0])

	return &models.User{
		ID:      int(value.ID), //nolint:gosec
		Login:   value.Login,
		Name:    value.Name,
		Blocked: value.Blocked,
	}, nil
}

func (s *Users) GetUsers(ctx context.Context) ([]*models.User, error) {
	const op = "DbUsers.GetUsers"

	resp, err := s.c.Connection.Select(space, "pk", 0, limit, tarantool.IterAll, tarantoolclient.Tuple{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users | %s:%w", op, err)
	}

	users := make([]*models.User, 0)

	for _, value := range resp.Tuples() {
		u := s.tupleToUser(value)

		users = append(users, &models.User{
			ID:      int(u.ID), //nolint:gosec
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

		if _, err := s.c.Connection.Insert(
			space,
			s.userCreatedToTuple(
				tarantoolclient.UserCreated{
					Name:    user.Name,
					Login:   user.Login,
					Blocked: user.Blocked,
				},
			),
		); err != nil {
			return nil, fmt.Errorf("failed to create user | %s:%w", op, err)
		}

		return s.GetUser(ctx, user.Login)
	}

	return nil, fmt.Errorf("failed to create user | %s:%w", op, dberrors.ErrUserAlreadyExists)
}

func (s *Users) UpdateUser(ctx context.Context, user models.UserUpdated) (*models.User, error) {
	const op = "DbUsers.UpdateUser"

	u, err := s.GetUser(ctx, user.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to update user | %s:%w", op, err)
	}

	if _, err := s.c.Connection.Replace(
		space,
		s.userUpdatedToTuple(
			tarantoolclient.UserUpdated{
				ID:      uint64(u.ID), //nolint:gosec
				Name:    user.Name,
				Login:   u.Login,
				Blocked: user.Blocked,
			},
		),
	); err != nil {
		return nil, fmt.Errorf("failed to update user | %s:%w", op, err)
	}

	return s.GetUser(ctx, user.Login)
}

func (s *Users) DeleteUser(ctx context.Context, login string) error {
	const op = "DbUsers.DeleteUser"

	if _, err := s.GetUser(ctx, login); err != nil {
		return fmt.Errorf("failed to delete user | %s:%w", op, err)
	}

	if _, err := s.c.Connection.Delete(space, "secondary", tarantoolclient.Tuple{login}); err != nil {
		return fmt.Errorf("failed to delete user | %s:%w", op, err)
	}

	return nil
}

func (s *Users) tupleToUser(tuple tarantoolclient.Tuple) tarantoolclient.User {
	return tarantoolclient.User{
		ID:      tuple[0].(uint64), //nolint:forcetypeassert
		Name:    tuple[1].(string), //nolint:forcetypeassert
		Login:   tuple[2].(string), //nolint:forcetypeassert
		Blocked: tuple[3].(bool),   //nolint:forcetypeassert
	}
}

func (s *Users) userCreatedToTuple(user tarantoolclient.UserCreated) tarantoolclient.Tuple {
	return tarantoolclient.Tuple{
		nil,
		user.Name,
		user.Login,
		user.Blocked,
	}
}

func (s *Users) userUpdatedToTuple(user tarantoolclient.UserUpdated) tarantoolclient.Tuple {
	return tarantoolclient.Tuple{
		user.ID,
		user.Name,
		user.Login,
		user.Blocked,
	}
}
