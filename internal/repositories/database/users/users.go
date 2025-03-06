package dbusers

import (
	"context"
	"fmt"

	dbclient "github.com/vtievsky/auth-id/internal/repositories/database/client"
)

type Client interface {
	GetUsers(ctx context.Context) ([]*dbclient.User, error)
	CreateUser(ctx context.Context, user dbclient.UserCreated) (*dbclient.User, error)
}

type UsersOpts struct {
	Client Client
}

type Users struct {
	client Client
}

func New(opts *UsersOpts) *Users {
	return &Users{
		client: opts.Client,
	}
}

func (s *Users) GetUsers(ctx context.Context) ([]*User, error) {
	const op = "DBUsers.GetUsers"

	ul, err := s.client.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users | %s:%w", op, err)
	}

	users := make([]*User, 0, len(ul))

	for _, user := range ul {
		users = append(users, &User{
			ID:       user.ID,
			Login:    user.Login,
			FullName: user.FullName,
			Blocked:  user.Blocked,
		})
	}

	return users, nil
}

func (s *Users) CreateUser(ctx context.Context, user UserCreated) (*User, error) {
	const op = "DBUsers.CreateUser"

	u, err := s.client.CreateUser(ctx, dbclient.UserCreated{
		Login:    user.Login,
		FullName: user.FullName,
		Blocked:  user.Blocked,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user | %s:%w", op, err)
	}

	return &User{
		ID:       u.ID,
		Login:    u.Login,
		FullName: u.FullName,
		Blocked:  u.Blocked,
	}, nil
}
