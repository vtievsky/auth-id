package sessionsvc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	userprivilegesvc "github.com/vtievsky/auth-id/internal/services/user-privileges"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
	"go.uber.org/zap"
)

type UserSvc interface {
	GetUser(ctx context.Context, login string) (*usersvc.User, error)
	ComparePassword(password, current []byte) error
}

type UserPrivilegeSvc interface {
	GetUserPrivileges(ctx context.Context, login string) ([]*userprivilegesvc.UserPrivilege, error)
}

type SessionSvcOpts struct {
	Logger           *zap.Logger
	UserSvc          UserSvc
	UserPrivilegeSvc UserPrivilegeSvc
}

type SessionSvc struct {
	logger           *zap.Logger
	userSvc          UserSvc
	userPrivilegeSvc UserPrivilegeSvc
}

func New(opts *SessionSvcOpts) *SessionSvc {
	return &SessionSvc{
		logger:           opts.Logger,
		userSvc:          opts.UserSvc,
		userPrivilegeSvc: opts.UserPrivilegeSvc,
	}
}

func (s *SessionSvc) Login(ctx context.Context, login, password string) ([]byte, error) {
	const op = "SessionSvc.Login"

	u, err := s.userSvc.GetUser(ctx, login)
	if err != nil {
		s.logger.Error("failed to get user",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	// Проверка пароля
	if err = s.userSvc.ComparePassword([]byte(u.Password), []byte(password)); err != nil {
		s.logger.Error("failed to compare password",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to compare password | %s:%w", op, err)
	}

	// Получение привилегий пользователя и создание сессии
	privileges, err := s.userPrivilegeSvc.GetUserPrivileges(ctx, u.Login)
	if err != nil {
		s.logger.Error("failed to fetch user privileges",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to fetch user privileges | %s:%w", op, err)
	}

	sessionID := uuid.NewString()

	/***/
	fmt.Printf("privileges of %s (%s)\n", login, sessionID)

	for _, privilege := range privileges {
		fmt.Printf("%s\n", privilege.Code)
	}
	/***/

	return []byte(sessionID), nil
}
