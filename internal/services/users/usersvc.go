package usersvc

import (
	"context"

	"go.uber.org/zap"
)

type User struct {
	ID        int
	Login     string
	FullName  string
	IsBlocked bool
}

type UserCreated struct {
	Login     string
	FullName  string
	IsBlocked bool
}

type UserUpdated struct {
	Login     string
	FullName  string
	IsBlocked bool
}

type UserSvcOpts struct {
	Logger *zap.Logger
}

type UserSvc struct {
	logger *zap.Logger
}

func New(opts *UserSvcOpts) *UserSvc {
	return &UserSvc{
		logger: opts.Logger,
	}
}

func (s *UserSvc) User(ctx context.Context, login string) (*User, error) {
	return &User{
		ID:        0,
		Login:     "",
		FullName:  "",
		IsBlocked: false,
	}, nil
}

func (s *UserSvc) GetUsers(ctx context.Context) ([]*User, error) {
	ul := make([]*User, 0)

	return ul, nil
}

func (s *UserSvc) CreateUser(ctx context.Context, user UserCreated) (*User, error) {
	// TODO Создание пользователя

	u, err := s.User(ctx, user.Login)
	if err != nil {

	}

	return u, nil
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
