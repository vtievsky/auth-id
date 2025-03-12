package usersvc

import (
	"context"
	"fmt"
	"strings"

	"github.com/vtievsky/auth-id/internal/repositories/models"
	"github.com/vtievsky/auth-id/pkg/cache"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint64
	Name     string
	Login    string
	Password string
	Blocked  bool
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

type UserUpdatedWithPass struct {
	Name     string
	Login    string
	Password string
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
			ID:       user.ID,
			Name:     user.Name,
			Login:    user.Login,
			Password: user.Password,
			Blocked:  user.Blocked,
		})
	}

	return users, nil
}

func (s *UserSvc) CreateUser(ctx context.Context, user UserCreated) (*User, error) {
	const op = "UserSvc.CreateUser"

	if strings.TrimSpace(user.Name) == "" {
		s.logger.Error("failed to create user",
			zap.String("login", user.Login),
			zap.Error(ErrInvalidName),
		)

		return nil, fmt.Errorf("failed to create user | %s:%w", op, ErrInvalidName)
	}

	if strings.TrimSpace(user.Login) == "" {
		s.logger.Error("failed to create user",
			zap.String("login", user.Login),
			zap.Error(ErrInvalidLogin),
		)

		return nil, fmt.Errorf("failed to create user | %s:%w", op, ErrInvalidLogin)
	}

	if strings.TrimSpace(user.Password) == "" {
		s.logger.Error("failed to create user",
			zap.String("login", user.Login),
			zap.Error(ErrInvalidPassword),
		)

		return nil, fmt.Errorf("failed to create user | %s:%w", op, ErrInvalidPassword)
	}

	hash, err := s.generateHashPassword([]byte(user.Password))
	if err != nil {
		s.logger.Error("failed to create user",
			zap.String("login", user.Login),
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
			zap.String("login", user.Login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to create user | %s:%w", op, err)
	}

	s.cacheByID.Add(u.ID, u)
	s.cacheByLogin.Add(u.Login, u)

	return &User{
		ID:       u.ID,
		Name:     u.Name,
		Login:    u.Login,
		Password: u.Password,
		Blocked:  u.Blocked,
	}, nil
}

// При "обычном" изменении пользователя пароль не изменяется
func (s *UserSvc) UpdateUser(ctx context.Context, user UserUpdated) (*User, error) {
	const op = "UserSvc.UpdateUser"

	u, err := s.GetUserByLogin(ctx, user.Login)
	if err != nil {
		s.logger.Error("failed to get user",
			zap.String("login", user.Login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	userUpdated, err := s.updateUser(ctx, UserUpdatedWithPass{
		Name:     user.Name,
		Login:    user.Login,
		Password: u.Password,
		Blocked:  user.Blocked,
	})
	if err != nil {
		s.logger.Error("failed to update user",
			zap.String("login", user.Login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to update user | %s:%w", op, err)
	}

	return userUpdated, nil
}

// Смена пароля пользователя с предварительной проверкой
func (s *UserSvc) ChangePass(ctx context.Context, login, current, changed string) error {
	const op = "UserSvc.ChangePass"

	u, err := s.GetUserByLogin(ctx, login)
	if err != nil {
		s.logger.Error("failed to get user",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	// Проверка текущего пароля
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(current)); err != nil {
		s.logger.Error("failed to compare password",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to compare password | %s:%w", op, err)
	}

	hash, err := s.generateHashPassword([]byte(changed))
	if err != nil {
		s.logger.Error("failed to generate hash password",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to generate hash password | %s:%w", op, err)
	}

	if _, err = s.updateUser(ctx, UserUpdatedWithPass{
		Name:     u.Name,
		Login:    u.Login,
		Password: string(hash),
		Blocked:  u.Blocked,
	}); err != nil {
		s.logger.Error("failed to update password",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to update password | %s:%w", op, err)
	}

	// Удалим старые данные пользователя из кеша
	s.cacheByID.Del(u.ID)
	s.cacheByLogin.Del(u.Login)

	return nil
}

// Смена пароля пользователя без предварительной проверки
func (s *UserSvc) ResetPass(ctx context.Context, login, changed string) error {
	const op = "UserSvc.ResetPass"

	u, err := s.GetUserByLogin(ctx, login)
	if err != nil {
		s.logger.Error("failed to get user",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	hash, err := s.generateHashPassword([]byte(changed))
	if err != nil {
		s.logger.Error("failed to generate hash password",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to generate hash password | %s:%w", op, err)
	}

	if _, err = s.updateUser(ctx, UserUpdatedWithPass{
		Name:     u.Name,
		Login:    u.Login,
		Password: string(hash),
		Blocked:  u.Blocked,
	}); err != nil {
		s.logger.Error("failed to update password",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to update password | %s:%w", op, err)
	}

	// Удалим старые данные пользователя из кеша
	s.cacheByID.Del(u.ID)
	s.cacheByLogin.Del(u.Login)

	return nil
}

func (s *UserSvc) DeleteUser(ctx context.Context, login string) error {
	const op = "UserSvc.DeleteUser"

	if err := s.storage.DeleteUser(ctx, login); err != nil {
		s.logger.Error("failed to delete user",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to delete user | %s:%w", op, err)
	}

	return nil
}

func (s *UserSvc) updateUser(ctx context.Context, user UserUpdatedWithPass) (*User, error) {
	if strings.TrimSpace(user.Name) == "" {
		return nil, ErrInvalidName
	}

	if strings.TrimSpace(user.Login) == "" {
		return nil, ErrInvalidLogin
	}

	if strings.TrimSpace(user.Password) == "" {
		return nil, ErrInvalidPassword
	}

	u, err := s.storage.UpdateUser(ctx, models.UserUpdated{
		Name:     user.Name,
		Login:    user.Login,
		Password: user.Password,
		Blocked:  user.Blocked,
	})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &User{
		ID:       u.ID,
		Name:     u.Name,
		Login:    u.Login,
		Password: u.Password,
		Blocked:  u.Blocked,
	}, nil
}

func (s *UserSvc) generateHashPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrGeneratePassword, err)
	}

	return hash, nil
}
