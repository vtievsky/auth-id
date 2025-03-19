package sessionsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	reposessions "github.com/vtievsky/auth-id/internal/repositories/sessions/sessions"
	userprivilegesvc "github.com/vtievsky/auth-id/internal/services/user-privileges"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
	authidjwt "github.com/vtievsky/auth-id/pkg/jwt"
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
	Store(ctx context.Context, login, sessionID string, privileges []string, ttl time.Duration) error
	Delete(ctx context.Context, login, sessionID string) error
	Find(ctx context.Context, sessionID, privilegeCode string) error
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
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	SigningKey       string
}

type SessionSvc struct {
	logger           *zap.Logger
	storage          Storage
	userSvc          UserSvc
	userPrivilegeSvc UserPrivilegeSvc
	sessionTTL       time.Duration
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	signingKey       string
}

func New(opts *SessionSvcOpts) *SessionSvc {
	return &SessionSvc{
		logger:           opts.Logger,
		storage:          opts.Storage,
		userSvc:          opts.UserSvc,
		userPrivilegeSvc: opts.UserPrivilegeSvc,
		accessTokenTTL:   opts.AccessTokenTTL,
		refreshTokenTTL:  opts.RefreshTokenTTL,
		sessionTTL:       opts.SessionTTL,
		signingKey:       opts.SigningKey,
	}
}

func (s *SessionSvc) Find(ctx context.Context, sessionID, privilegeCode string) error {
	const op = "SessionSvc.Find"

	if err := s.storage.Find(ctx, sessionID, privilegeCode); err != nil {
		s.logger.Error("failed to search session privilege",
			zap.String("session_id", sessionID),
			zap.String("privilege_code", privilegeCode),
			zap.Error(err),
		)

		return fmt.Errorf("failed to search session privilege | %s:%w", op, err)
	}

	return nil
}

func (s *SessionSvc) Login(ctx context.Context, login, password string) (*Tokens, error) {
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

	var sessionPrivilegesExpiredAt time.Time // Дата окончания срока действия всех привилегий

	for _, privilege := range privileges {
		sessionPrivileges = append(sessionPrivileges, privilege.Code)

		if sessionPrivilegesExpiredAt.Before(privilege.DateOut) {
			sessionPrivilegesExpiredAt = privilege.DateOut
		}
	}

	tokens, err := s.generateTokens(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to generate tokens",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to generate tokens | %s:%w", op, err)
	}

	// Сохранение сессии и ее привилегий
	sessionDuration := time.Until(sessionPrivilegesExpiredAt)
	sessionDuration = s.compareSessionWithPrivilegesTTL(s.sessionTTL, sessionDuration)

	// Длительность хранения списка привилегий не может быть дольше
	// общей длительности сессии пользователя (refreshTokenTTL)
	sessionDuration = s.compareSessionWithRefreshTokenTTL(sessionDuration, s.refreshTokenTTL)

	if err = s.storage.Store(ctx, login, sessionID, sessionPrivileges, sessionDuration); err != nil {
		s.logger.Error("failed to store session",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to store session | %s:%w", op, err)
	}

	s.logger.Debug("Registration was successful",
		zap.String("login", login),
		zap.String("session_id", sessionID),
	)

	return tokens, nil
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

func (s *SessionSvc) Delete(ctx context.Context, login, sessionID string) error {
	const op = "SessionSvc.Delete"

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

func (s *SessionSvc) generateTokens(_ context.Context, sessionID string) (*Tokens, error) {
	const op = "SessionSvc.generateTokens"

	if s.refreshTokenTTL < s.accessTokenTTL {
		return nil, ErrInvalidAccessTokenTTL
	}

	var (
		accessToken  []byte
		refreshToken []byte
		signingKey   = []byte(s.signingKey)
		current      = time.Now()
	)

	g := errgroup.Group{}

	g.Go(func() error {
		var err error

		accessToken, err = authidjwt.NewAccessToken(signingKey, &authidjwt.TokenOpts{
			SessionID: sessionID,
			ExpiredAt: current.Add(s.accessTokenTTL),
		})
		if err != nil {
			return fmt.Errorf("failed to generate access token | %s:%w", op, err)
		}

		return nil
	})

	g.Go(func() error {
		var err error

		refreshToken, err = authidjwt.NewRefreshToken(signingKey, &authidjwt.TokenOpts{
			SessionID: sessionID,
			ExpiredAt: current.Add(s.refreshTokenTTL),
		})
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

func (s *SessionSvc) compareSessionWithPrivilegesTTL(sessionTTL, sessionPrivilegesTTL time.Duration) time.Duration {
	if sessionTTL < sessionPrivilegesTTL {
		s.logger.Debug("the duration of the session is less than the duration of the privileges. "+
			"The duration of the session will be used",
			zap.String("session_ttl", sessionTTL.String()),
			zap.String("session_privileges_ttl", sessionPrivilegesTTL.String()),
		)

		return sessionTTL
	}

	return sessionPrivilegesTTL
}

func (s *SessionSvc) compareSessionWithRefreshTokenTTL(sessionTTL, refreshTokenTTL time.Duration) time.Duration {
	if sessionTTL < refreshTokenTTL {
		s.logger.Debug("the duration of the session is less than the duration of the refresh token. "+
			"The duration of the session will be used",
			zap.String("session_ttl", sessionTTL.String()),
			zap.String("refresh_token_ttl", refreshTokenTTL.String()),
		)

		return sessionTTL
	}

	return refreshTokenTTL
}
