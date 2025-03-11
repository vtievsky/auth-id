package usersvc

import (
	"context"
	"fmt"

	"github.com/vtievsky/auth-id/internal/repositories/models"
	"github.com/vtievsky/auth-id/pkg/cache"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID      uint64
	Name    string
	Login   string
	Blocked bool
}

type UserCreated struct {
	Name     string
	Login    string
	Password string
	Blocked  bool
}

type UserUpdated struct {
	Name    string
	Login   string
	Blocked bool
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
	logger       *zap.Logger
	storage      Storage
	cacheByID    cache.Cache[uint64, *models.User]
	cacheByLogin cache.Cache[string, *models.User]
}

func New(opts *UserSvcOpts) *UserSvc {
	return &UserSvc{
		logger:       opts.Logger,
		storage:      opts.Storage,
		cacheByID:    cache.New[uint64, *models.User](),
		cacheByLogin: cache.New[string, *models.User](),
	}
}

func (s *UserSvc) GetUser(ctx context.Context, login string) (*User, error) {
	return s.GetUserByLogin(ctx, login)
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
			ID:      user.ID,
			Name:    user.Name,
			Login:   user.Login,
			Blocked: user.Blocked,
		})
	}

	return users, nil
}

func (s *UserSvc) CreateUser(ctx context.Context, user UserCreated) (*User, error) {
	const op = "UserSvc.CreateUser"

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to generate has from password",
			zap.String("username", user.Login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to create user | %s:%w", op, err)
	}

	u, err := s.storage.CreateUser(ctx, models.UserCreated{
		Name:     user.Name,
		Login:    user.Login,
		Password: string(hash),
		Blocked:  user.Blocked,
	})
	if err != nil {
		s.logger.Error("failed to create user",
			zap.String("username", user.Login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to create user | %s:%w", op, err)
	}

	s.cacheByID.Add(u.ID, u)
	s.cacheByLogin.Add(u.Login, u)

	return &User{
		ID:      u.ID,
		Name:    u.Name,
		Login:   u.Login,
		Blocked: u.Blocked,
	}, nil
}

func (s *UserSvc) UpdateUser(ctx context.Context, user UserUpdated) (*User, error) {
	const op = "UserSvc.UpdateUser"

	u, err := s.storage.UpdateUser(ctx, models.UserUpdated{
		Name:    user.Name,
		Login:   user.Login,
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
		ID:      u.ID,
		Name:    u.Name,
		Login:   u.Login,
		Blocked: u.Blocked,
	}, nil
}

func (s *UserSvc) DeleteUser(ctx context.Context, login string) error {
	const op = "UserSvc.DeleteUser"

	if err := s.storage.DeleteUser(ctx, login); err != nil {
		s.logger.Error("failed to delete user",
			zap.String("username", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to delete user | %s:%w", op, err)
	}

	return nil
}
