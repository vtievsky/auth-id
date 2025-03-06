package usersvc

import (
	"context"
	"fmt"

	dbusers "github.com/vtievsky/auth-id/internal/repositories/database/users"
	"go.uber.org/zap"
)

type User struct {
	ID       int
	Login    string
	FullName string
	Blocked  bool
}

type UserCreated struct {
	Login    string
	FullName string
	Blocked  bool
}

type UserUpdated struct {
	Login    string
	FullName string
	Blocked  bool
}

type Storage interface {
	GetUsers(ctx context.Context) ([]*dbusers.User, error)
	CreateUser(ctx context.Context, user dbusers.UserCreated) (*dbusers.User, error)
}

type UserSvcOpts struct {
	Logger  *zap.Logger
	Storage Storage
}

type UserSvc struct {
	logger  *zap.Logger
	storage Storage
}

func New(opts *UserSvcOpts) *UserSvc {
	return &UserSvc{
		logger:  opts.Logger,
		storage: opts.Storage,
	}
}

func (s *UserSvc) User(ctx context.Context, login string) (*User, error) {
	return &User{
		ID:       0,
		Login:    "",
		FullName: "",
		Blocked:  false,
	}, nil
}

func (s *UserSvc) GetUsers(ctx context.Context) ([]*User, error) {
	const op = "UserSvc.GetUsers"

	ul, err := s.storage.GetUsers(ctx)
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

func (s *UserSvc) CreateUser(ctx context.Context, user UserCreated) (*User, error) {
	const op = "UserSvc.CreateUser"

	u, err := s.storage.CreateUser(ctx, dbusers.UserCreated{
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

func (s *UserSvc) UpdateUser(ctx context.Context, user UserUpdated) (*User, error) {
	// TODO Изменение пользователя

	u, err := s.User(ctx, user.Login)
	if err != nil {

	}

	return u, nil
}

func (s *UserSvc) DeleteUser(ctx context.Context, login string) error {
	// TODO Изменение пользователя

	return nil
}
