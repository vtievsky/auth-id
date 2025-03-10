package usersvc

import (
	"context"
	"fmt"
	"sync"
	"time"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	"go.uber.org/zap"
)

const (
	cacheTTL = time.Second * 60
)

type User struct {
	ID      int
	Name    string
	Login   string
	Blocked bool
}

type UserCreated struct {
	Name    string
	Login   string
	Blocked bool
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
	lastTime     time.Time
	cacheByID    map[int]*models.User
	cacheByLogin map[string]*models.User
	mu           sync.RWMutex
}

func New(opts *UserSvcOpts) *UserSvc {
	return &UserSvc{
		logger:       opts.Logger,
		storage:      opts.Storage,
		lastTime:     time.Time{},
		cacheByID:    make(map[int]*models.User),
		cacheByLogin: make(map[string]*models.User),
		mu:           sync.RWMutex{},
	}
}

func (s *UserSvc) GetUser(ctx context.Context, login string) (*User, error) {
	return s.GetUserByLogin(ctx, login)
}

func (s *UserSvc) GetUserByID(ctx context.Context, id int) (*User, error) {
	const op = "UserSvc.GetUserByID"

	s.mu.RLock()

	// Кеш невалидный
	if cacheTTL < time.Since(s.lastTime) {
		s.mu.RLocker().Unlock()

		if err := s.syncUsers(ctx); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		s.mu.RLock()
	}

	defer s.mu.RLocker().Unlock()

	if val, ok := s.cacheByID[id]; ok {
		return &User{
			ID:      val.ID,
			Login:   val.Login,
			Name:    val.Name,
			Blocked: val.Blocked,
		}, nil
	}

	return nil, dberrors.ErrUserNotFound
}

func (s *UserSvc) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	const op = "UserSvc.GetUserByLogin"

	s.mu.RLock()

	// Кеш невалидный
	if cacheTTL < time.Since(s.lastTime) {
		s.mu.RLocker().Unlock()

		if err := s.syncUsers(ctx); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		s.mu.RLock()
	}

	defer s.mu.RLocker().Unlock()

	if val, ok := s.cacheByLogin[login]; ok {
		return &User{
			ID:      val.ID,
			Login:   val.Login,
			Name:    val.Name,
			Blocked: val.Blocked,
		}, nil
	}

	return nil, dberrors.ErrUserNotFound
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
			Login:   user.Login,
			Name:    user.Name,
			Blocked: user.Blocked,
		})
	}

	return users, nil
}

func (s *UserSvc) CreateUser(ctx context.Context, user UserCreated) (*User, error) {
	const op = "UserSvc.CreateUser"

	u, err := s.storage.CreateUser(ctx, models.UserCreated{
		Login:   user.Login,
		Name:    user.Name,
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
		ID:      u.ID,
		Login:   u.Login,
		Name:    u.Name,
		Blocked: u.Blocked,
	}, nil
}

func (s *UserSvc) UpdateUser(ctx context.Context, user UserUpdated) (*User, error) {
	const op = "UserSvc.UpdateUser"

	u, err := s.storage.UpdateUser(ctx, models.UserUpdated{
		Login:   user.Login,
		Name:    user.Name,
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
		Login:   u.Login,
		Name:    u.Name,
		Blocked: u.Blocked,
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

func (s *UserSvc) syncUsers(ctx context.Context) error {
	const op = "UserSvc.syncUsers"

	s.mu.Lock()
	defer s.mu.Unlock()

	resp, err := s.storage.GetUsers(ctx)
	if err != nil {
		s.logger.Error("failed to sync users",
			zap.Error(err),
		)

		return fmt.Errorf("failed to sync users | %s:%w", op, err)
	}

	for _, user := range resp {
		s.cacheByID[user.ID] = user
		s.cacheByLogin[user.Login] = user
	}

	// Зафиксируем время синхронизации справочника
	s.lastTime = time.Now()

	return nil
}
