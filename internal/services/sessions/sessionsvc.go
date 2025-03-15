package sessionsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	redissessions "github.com/vtievsky/auth-id/internal/repositories/sessions/sessions"
	userprivilegesvc "github.com/vtievsky/auth-id/internal/services/user-privileges"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
	"go.uber.org/zap"
)

type Session struct {
	ID        []byte
	CreatedAt time.Time
	ExpiredAt time.Time
}

type Storage interface {
	List(ctx context.Context, login string) ([]*redissessions.Session, error)
	Find(ctx context.Context, sessionID, privilege string) error
	Store(ctx context.Context, login, sessionID string, privileges []string, ttl time.Duration) error
	Delete(ctx context.Context, login, sessionID string) error
}

type UserSvc interface {
	GetUser(ctx context.Context, login string) (*usersvc.User, error)
	ComparePassword(password, current []byte) error
}

type UserPrivilegeSvc interface {
	GetUserPrivileges(ctx context.Context, login string) ([]*userprivilegesvc.UserPrivilege, error)
}

type SessionSvcOpts struct {
	Logger           *zap.Logger
	Storage          Storage
	UserSvc          UserSvc
	UserPrivilegeSvc UserPrivilegeSvc
	SessionTTL       time.Duration
}

type SessionSvc struct {
	logger           *zap.Logger
	storage          Storage
	userSvc          UserSvc
	userPrivilegeSvc UserPrivilegeSvc
	sessionTTL       time.Duration
}

func New(opts *SessionSvcOpts) *SessionSvc {
	return &SessionSvc{
		logger:           opts.Logger,
		storage:          opts.Storage,
		userSvc:          opts.UserSvc,
		userPrivilegeSvc: opts.UserPrivilegeSvc,
		sessionTTL:       opts.SessionTTL,
	}
}

func (s *SessionSvc) Login(ctx context.Context, login, password string) (*Session, error) {
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
	sessionPrivileges := make([]string, 0, len(privileges))

	var expiredAt time.Time // Время окончания действия всех привилегий

	for _, privilege := range privileges {
		sessionPrivileges = append(sessionPrivileges, privilege.Code)

		if expiredAt.Before(privilege.DateOut) {
			expiredAt = privilege.DateOut
		}
	}

	// Сохранение сессии и ее привилегий
	current := time.Now()
	sessionDuration := time.Until(expiredAt)
	sessionDuration = min(s.sessionTTL, sessionDuration)

	if err = s.storage.Store(ctx, login, sessionID, sessionPrivileges, sessionDuration); err != nil {
		s.logger.Error("failed to store session",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to store session | %s:%w", op, err)
	}

	return &Session{
		ID:        []byte(sessionID),
		CreatedAt: current,
		ExpiredAt: expiredAt,
	}, nil
}

func (s *SessionSvc) GetUserSessions(ctx context.Context, login string) ([]*Session, error) {
	const op = "SessionSvc.GetUserSessions"

	u, err := s.userSvc.GetUser(ctx, login)
	if err != nil {
		s.logger.Error("failed to get user",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	sessions, err := s.storage.List(ctx, u.Login)
	if err != nil {
		s.logger.Error("failed to get user sessions",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get user sessions | %s:%w", op, err)
	}

	var (
		current   = time.Now()
		createdAt time.Time
		expiredAt time.Time
	)

	ul := make([]*Session, 0, len(sessions))

	for _, session := range sessions {
		expiredAt = current.Add(session.TTL)
		createdAt = expiredAt.Add(-s.sessionTTL)

		ul = append(ul, &Session{
			ID:        []byte(session.ID),
			CreatedAt: createdAt,
			ExpiredAt: expiredAt,
		})
	}

	return ul, nil
}

func (s *SessionSvc) Delete(ctx context.Context, login, sessionID string) error {
	const op = "SessionSvc.Delete"

	if err := s.storage.Delete(ctx, login, sessionID); err != nil {
		s.logger.Error("failed to delete session",
			zap.String("login", login),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)

		return fmt.Errorf("failed to delete session | %s:%w", op, err)
	}

	return nil
}
