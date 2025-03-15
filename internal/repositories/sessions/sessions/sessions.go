package reposessions

import (
	"context"
	"fmt"
	"strings"
	"time"

	clientredis "github.com/vtievsky/auth-id/internal/repositories/sessions/client/redis"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	limit = 5
	space = "omo:"
)

type sessionStats struct {
	ID  string
	TTL time.Duration
}

type Session struct {
	ID  string
	TTL time.Duration
}

type SessionsOpts struct {
	Logger *zap.Logger
	Client *clientredis.Client
}

type Sessions struct {
	logger *zap.Logger
	client *clientredis.Client
}

func New(opts *SessionsOpts) *Sessions {
	return &Sessions{
		logger: opts.Logger,
		client: opts.Client,
	}
}

func (s *Sessions) Find(ctx context.Context, sessionID, privilege string) error {
	const op = "Sessions.Find"

	key := s.keySession(sessionID)
	command := s.client.SIsMember(ctx, key, privilege)

	switch {
	case command.Err() != nil:
		return fmt.Errorf("failed to search session privilege | %s:%w", op, command.Err())
	case command.Val():
		return nil
	}

	return fmt.Errorf("%s:%w", op, ErrSessionPrivilegeNotFound)
}

func (s *Sessions) List(ctx context.Context, login string) ([]*Session, error) {
	const op = "Sessions.List"

	fetchTTL := func(actx context.Context, acombine chan<- sessionStats, asessionID string) func() error {
		return func() error {
			ttl, err := s.client.TTL(actx, asessionID).Result()
			if err != nil {
				s.logger.Error("failed to get session ttl",
					zap.String("session_id", asessionID),
					zap.Error(err),
				)

				return nil
			}

			acombine <- sessionStats{
				ID:  asessionID,
				TTL: ttl,
			}

			return nil
		}
	}

	ul, err := s.client.SMembers(ctx, s.keySessions(login)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions list | %s:%w", op, err)
	}

	combine := make(chan sessionStats)

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(limit + 1)

	g.Go(func() error {
		for _, session := range ul {
			g.Go(fetchTTL(gCtx, combine, session))
		}

		return nil
	})

	go func() {
		err = g.Wait()

		close(combine)
	}()

	sessions := make([]*Session, 0)

	for v := range combine {
		sessions = append(sessions, &Session{
			ID:  s.sanitizeID(v.ID),
			TTL: v.TTL,
		})
	}

	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return sessions, nil
}

func (s *Sessions) sanitizeID(sessionID string) string {
	return strings.Replace(sessionID, space, "", 1)
}

// Ключ сессии (с множеством привилегий)
func (s *Sessions) keySession(sessionID string) string {
	return fmt.Sprintf("%s%s", space, sessionID)
}

// Ключ с сессиями пользователя
func (s *Sessions) keySessions(login string) string {
	return fmt.Sprintf("%s%s", space, login)
}
