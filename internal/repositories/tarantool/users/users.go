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
	space         = "user"
	spaceUserRole = "role_user"
	limit         = 25
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

	user := s.tupleToUser(resp.Tuples()[0])

	return &models.User{
		ID:       user.ID,
		Name:     user.Name,
		Login:    user.Login,
		Password: user.Password,
		Blocked:  user.Blocked,
	}, nil
}

func (s *Users) GetUsers(ctx context.Context) ([]*models.User, error) {
	const op = "DbUsers.GetUsers"

	resp, err := s.c.Connection.Select(space, "pk", 0, limit, tarantool.IterAll, tarantoolclient.Tuple{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users | %s:%w", op, err)
	}

	var user tarantoolclient.User

	users := make([]*models.User, 0, len(resp.Tuples()))

	for _, tuple := range resp.Tuples() {
		user = s.tupleToUser(tuple)

		users = append(users, &models.User{
			ID:       user.ID,
			Name:     user.Name,
			Login:    user.Login,
			Password: user.Password,
			Blocked:  user.Blocked,
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

		userCreated := tarantoolclient.UserCreated{
			Name:     user.Name,
			Login:    user.Login,
			Password: user.Password,
			Blocked:  user.Blocked,
		}

		if _, err := s.c.Connection.Insert(space, userCreated.ToTuple()); err != nil {
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

	userUpdated := tarantoolclient.UserUpdated{
		ID:       u.ID,
		Name:     user.Name,
		Login:    u.Login,
		Password: u.Password,
		Blocked:  user.Blocked,
	}

	if _, err := s.c.Connection.Replace(space, userUpdated.ToTuple()); err != nil {
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
		ID:       tuple[0].(uint64), //nolint:forcetypeassert
		Name:     tuple[1].(string), //nolint:forcetypeassert
		Login:    tuple[2].(string), //nolint:forcetypeassert
		Password: tuple[3].(string), //nolint:forcetypeassert
		Blocked:  tuple[4].(bool),   //nolint:forcetypeassert
	}
}
