package usersvc

import (
	"context"
	"fmt"

	"github.com/vtievsky/auth-id/internal/repositories/models"
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
	GetUser(ctx context.Context, login string) (*models.User, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
	CreateUser(ctx context.Context, user models.UserCreated) (*models.User, error)
	UpdateUser(ctx context.Context, user models.UserUpdated) (*models.User, error)
	DeleteUser(ctx context.Context, login string) error
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

func (s *UserSvc) GetUser(ctx context.Context, login string) (*User, error) {
	const op = "UserSvc.GetUser"

	resp, err := s.storage.GetUser(ctx, login)
	if err != nil {
		s.logger.Error("failed to get user",
			zap.String("username", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	return &User{
		ID:       resp.ID,
		Login:    resp.Login,
		FullName: resp.Name,
		Blocked:  resp.Blocked,
	}, nil
}

func (s *UserSvc) GetUsers(ctx context.Context) ([]*User, error) {
	const op = "UserSvc.GetUsers"

	ul, err := s.storage.GetUsers(ctx)
	if err != nil {
		s.logger.Error("failed to get users",
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get users | %s:%w", op, err)
	}

	users := make([]*User, 0, len(ul))

	for _, user := range ul {
		users = append(users, &User{
			ID:       user.ID,
			Login:    user.Login,
			FullName: user.Name,
			Blocked:  user.Blocked,
		})
	}

	return users, nil
}

func (s *UserSvc) CreateUser(ctx context.Context, user UserCreated) (*User, error) {
	const op = "UserSvc.CreateUser"

	u, err := s.storage.CreateUser(ctx, models.UserCreated{
		Login:   user.Login,
		Name:    user.FullName,
		Blocked: user.Blocked,
	})
	if err != nil {
		s.logger.Error("failed to create user",
			zap.String("username", user.Login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to create user | %s:%w", op, err)
	}

	return &User{
		ID:       u.ID,
		Login:    u.Login,
		FullName: u.Name,
		Blocked:  u.Blocked,
	}, nil
}

func (s *UserSvc) UpdateUser(ctx context.Context, user UserUpdated) (*User, error) {
	const op = "UserSvc.UpdateUser"

	u, err := s.storage.UpdateUser(ctx, models.UserUpdated{
		Login:   user.Login,
		Name:    user.FullName,
		Blocked: user.Blocked,
	})
	if err != nil {
		s.logger.Error("failed to update user",
			zap.String("username", user.Login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to update user | %s:%w", op, err)
	}

	return &User{
		ID:       u.ID,
		Login:    u.Login,
		FullName: u.Name,
		Blocked:  u.Blocked,
	}, nil
}

func (s *UserSvc) DeleteUser(ctx context.Context, login string) error {
	const op = "UserSvc.UpdateUser"

	if err := s.storage.DeleteUser(ctx, login); err != nil {
		s.logger.Error("failed to delete user",
			zap.String("username", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to delete user | %s:%w", op, err)
	}

	return nil
}
