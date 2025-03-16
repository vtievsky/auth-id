package authsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	reposessions "github.com/vtievsky/auth-id/internal/repositories/sessions/sessions"
	authtoken "github.com/vtievsky/auth-id/internal/services/auth/tokens"
	userprivilegesvc "github.com/vtievsky/auth-id/internal/services/user-privileges"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Session struct {
	ID        []byte
	CreatedAt time.Time
	ExpiredAt time.Time
}

type Storage interface {
	List(ctx context.Context, login string) ([]*reposessions.Session, error)
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

type AuthSvcOpts struct {
	Logger           *zap.Logger
	Storage          Storage
	UserSvc          UserSvc
	UserPrivilegeSvc UserPrivilegeSvc
	SessionTTL       time.Duration
	SigningKey       string
}

type AuthSvc struct {
	logger           *zap.Logger
	storage          Storage
	userSvc          UserSvc
	userPrivilegeSvc UserPrivilegeSvc
	sessionTTL       time.Duration
	signingKey       string
}

func New(opts *AuthSvcOpts) *AuthSvc {
	return &AuthSvc{
		logger:           opts.Logger,
		storage:          opts.Storage,
		userSvc:          opts.UserSvc,
		userPrivilegeSvc: opts.UserPrivilegeSvc,
		sessionTTL:       opts.SessionTTL,
		signingKey:       opts.SigningKey,
	}
}

func (s *AuthSvc) Login(ctx context.Context, login, password string) (*Tokens, error) {
	const op = "AuthSvc.Login"

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
	sessionDuration := time.Until(expiredAt)
	sessionDuration = min(s.sessionTTL, sessionDuration)

	if err = s.storage.Store(ctx, login, sessionID, sessionPrivileges, sessionDuration); err != nil {
		s.logger.Error("failed to store session",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to store session | %s:%w", op, err)
	}

	// Генерация токенов
	tokens, err := s.generateTokens(ctx, &authtoken.TokenOpts{
		SessionID: sessionID,
		ExpiredAt: expiredAt,
	})
	if err != nil {
		s.logger.Error("failed to generate tokens",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to generate tokens | %s:%w", op, err)
	}

	return tokens, nil
}

func (s *AuthSvc) GetUserSessions(ctx context.Context, login string) ([]*Session, error) {
	const op = "AuthSvc.GetUserSessions"

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
		expiredAt time.Time
	)

	ul := make([]*Session, 0, len(sessions))

	for _, session := range sessions {
		expiredAt = current.Add(session.TTL)

		ul = append(ul, &Session{
			ID:        []byte(session.ID),
			CreatedAt: session.CreatedAt,
			ExpiredAt: expiredAt,
		})
	}

	return ul, nil
}

func (s *AuthSvc) Delete(ctx context.Context, login, sessionID string) error {
	const op = "AuthSvc.Delete"

	u, err := s.userSvc.GetUser(ctx, login)
	if err != nil {
		s.logger.Error("failed to get user",
			zap.String("login", login),
			zap.Error(err),
		)

		return fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	if err := s.storage.Delete(ctx, u.Login, sessionID); err != nil {
		s.logger.Error("failed to delete session",
			zap.String("login", login),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)

		return fmt.Errorf("failed to delete session | %s:%w", op, err)
	}

	return nil
}

func (s *AuthSvc) generateTokens(_ context.Context, tokenOpts *authtoken.TokenOpts) (*Tokens, error) {
	const op = "AuthSvc.generateTokens"

	var (
		accessToken  []byte
		refreshToken []byte
		signingKey   = []byte(s.signingKey)
	)

	g := errgroup.Group{}

	g.Go(func() error {
		var err error

		accessToken, err = authtoken.NewAccessToken(signingKey, tokenOpts)
		if err != nil {
			return fmt.Errorf("failed to generate access token | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
		var err error

		refreshToken, err = authtoken.NewRefreshToken(signingKey, tokenOpts)
		if err != nil {
			return fmt.Errorf("failed to generate refresh token | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &Tokens{
		AccessToken:  string(accessToken),
		RefreshToken: string(refreshToken),
	}, nil
}
