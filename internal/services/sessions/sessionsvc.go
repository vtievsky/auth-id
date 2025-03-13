package sessionsvc

import (
	"context"
	"fmt"

	usersvc "github.com/vtievsky/auth-id/internal/services/users"
	"go.uber.org/zap"
)

type UserSvc interface {
	GetUser(ctx context.Context, login string) (*usersvc.User, error)
	ComparePassword(password, current []byte) error
}

type SessionSvcOpts struct {
	Logger  *zap.Logger
	UserSvc UserSvc
}

type SessionSvc struct {
	logger  *zap.Logger
	userSvc UserSvc
}

func New(opts *SessionSvcOpts) *SessionSvc {
	return &SessionSvc{
		logger:  opts.Logger,
		userSvc: opts.UserSvc,
	}
}

func (s *SessionSvc) Login(ctx context.Context, login, password string) error {
	const op = "SessionSvc.Login"

	u, err := s.userSvc.GetUser(ctx, login)
	if err != nil {
		s.logger.Error("failed to get user",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	// Проверка пароля
	if err = s.userSvc.ComparePassword([]byte(u.Password), []byte(password)); err != nil {
		s.logger.Error("failed to compare password",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to compare password | %s:%w", op, err)
	}

	// Получение привилегий пользователя и создание сессии

	return nil
}
