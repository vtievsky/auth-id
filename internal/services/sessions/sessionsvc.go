package sessionsvc

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	reposessions "github.com/vtievsky/auth-id/internal/repositories/sessions/sessions"
	userprivilegesvc "github.com/vtievsky/auth-id/internal/services/user-privileges"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
	"github.com/vtievsky/auth-id/pkg/cache"
	authidjwt "github.com/vtievsky/auth-id/pkg/jwt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var tracer = otel.Tracer("auth-id/sessions") //nolint:gochecknoglobals

const (
	MetricKindFailedGetUser         = "unknown_user"
	MetricKindInvalidPassword       = "invalid_password"
	MetricKindFailedFetchPrivileges = "failed_fetch_privileges"
	MetricKindEmptyPrivileges       = "empty_privileges"
	MetricKindFailedGenerateToken   = "failed_generate_token"
	MetricKindFailedStoreSession    = "failed_store_session"
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

type SessionCart struct {
	ID        string
	Login     string
	CreatedAt time.Time
}

type Storage interface {
	Get(ctx context.Context, sessionID string) (*reposessions.SessionCart, error)
	List(ctx context.Context, login string, pageSize, offset uint32) ([]*reposessions.Session, error)
	ListSessionPrivileges(ctx context.Context, sessionID string, pageSize, offset uint32) ([]string, error)
	Store(ctx context.Context, login, sessionID string, privileges []string, ttl time.Duration) error
	Delete(ctx context.Context, login, sessionID string) error
}

type UserSvc interface {
	GetUser(ctx context.Context, login string) (*usersvc.User, error)
	ComparePassword(password, current []byte) error
}

type UserPrivilegeSvc interface {
	GetUserPrivileges(ctx context.Context, login string, pageSize, offset uint32) ([]*userprivilegesvc.UserPrivilege, error)
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
	cacheByID        cache.Cache[string, []string]
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
		cacheByID:        cache.New[string, []string](),
	}
}

func (s *SessionSvc) Get(ctx context.Context, sessionID string) (*SessionCart, error) {
	const op = "SessionSvc.Get"

	val, err := s.storage.Get(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to get session cart",
			zap.String("session_id", sessionID),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get session cart | %s:%w", op, err)
	}

	return &SessionCart{
		ID:        val.ID,
		Login:     val.Login,
		CreatedAt: val.CreatedAt,
	}, nil
}

func (s *SessionSvc) Login(ctx context.Context, login, password string) (*Tokens, error) {
	const op = "SessionSvc.Login"

	ctx, span := tracer.Start(ctx, "login")
	defer span.End()

	span.AddEvent("start")

	u, err := s.userSvc.GetUser(ctx, login)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		incrLoginFail(ctx, MetricKindFailedGetUser)

		s.logger.Error("failed to get user",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	span.AddEvent("user has been received")

	// Проверка пароля
	if err = s.userSvc.ComparePassword([]byte(u.Password), []byte(password)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		incrLoginFail(ctx, MetricKindInvalidPassword)

		s.logger.Error("failed to compare password",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to compare password | %s:%w", op, err)
	}

	span.AddEvent("password has been compared")

	// Получение привилегий пользователя и создание сессии
	privileges, err := s.userPrivilegeSvc.GetUserPrivileges(ctx, u.Login, math.MaxUint32, 0)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		incrLoginFail(ctx, MetricKindFailedFetchPrivileges)

		s.logger.Error("failed to fetch user privileges",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to fetch user privileges | %s:%w", op, err)
	}

	if len(privileges) < 1 {
		err = ErrSessionPrivilegeNotFound

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		incrLoginFail(ctx, MetricKindEmptyPrivileges)

		s.logger.Error("empty user privileges",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("empty user privileges | %s:%w", op, err)
	}

	span.AddEvent("privileges has been received")

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
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		incrLoginFail(ctx, MetricKindFailedGenerateToken)

		s.logger.Error("failed to generate tokens",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to generate tokens | %s:%w", op, err)
	}

	span.AddEvent("token has been generated")

	// Сохранение сессии и ее привилегий
	sessionDuration := time.Until(sessionPrivilegesExpiredAt)
	sessionDuration = s.compareSessionWithPrivilegesTTL(s.sessionTTL, sessionDuration)

	// Длительность хранения списка привилегий не может быть дольше
	// общей длительности сессии пользователя (refreshTokenTTL)
	sessionDuration = s.compareSessionWithRefreshTokenTTL(sessionDuration, s.refreshTokenTTL)

	if err = s.storage.Store(ctx, login, sessionID, sessionPrivileges, sessionDuration); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		incrLoginFail(ctx, MetricKindFailedStoreSession)

		s.logger.Error("failed to store session",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to store session | %s:%w", op, err)
	}

	span.AddEvent("registration was successful")

	s.logger.Debug("Registration was successful",
		zap.String("login", login),
		zap.String("session_id", sessionID),
	)

	incrLoginSuccess(ctx)

	return tokens, nil
}

func (s *SessionSvc) GetUserSessions(ctx context.Context, login string, pageSize, offset uint32) ([]*Session, error) {
	const op = "SessionSvc.GetUserSessions"

	u, err := s.userSvc.GetUser(ctx, login)
	if err != nil {
		s.logger.Error("failed to get user",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get user | %s:%w", op, err)
	}

	sessions, err := s.storage.List(ctx, u.Login, pageSize, offset)
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
