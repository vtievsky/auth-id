package reposessions

import (
	"context"
	"fmt"
	"time"

	clientredis "github.com/vtievsky/auth-id/internal/repositories/sessions/client/redis"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	limit         = 5
	space         = "uev:"
	spaceCarts    = "pak:"
	spaceSessions = "omo:"
)

type sessionStats struct {
	ID        string
	TTL       time.Duration
	CreatedAt time.Time
}

type Session struct {
	ID        string
	TTL       time.Duration
	CreatedAt time.Time
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

func (s *Sessions) List(ctx context.Context, login string) ([]*Session, error) {
	const op = "Sessions.List"

	fetchStats := func(actx context.Context, acombine chan<- sessionStats, asessionID string) func() error {
		return func() error {
			var (
				ttl  time.Duration
				cart SessionCart
			)

			g, gCtx := errgroup.WithContext(actx)

			g.Go(func() error {
				var err error

				ttl, err = s.client.TTL(gCtx, s.keySession(asessionID)).Result()
				if err != nil {
					s.logger.Error("failed to get session ttl",
						zap.String("session_id", asessionID),
						zap.Error(err),
					)

					return nil
				}

				return nil
			})

			g.Go(func() error {
				cmd := s.client.HGetAll(gCtx, s.keyCart(asessionID))

				switch {
				case cmd.Err() != nil:
					s.logger.Error("failed to get session cart",
						zap.String("session_id", asessionID),
						zap.Error(cmd.Err()),
					)

					return nil //nolint:nilerr
				case len(cmd.Val()) == 0:
					// Скорее всего истек ttl ключа
					return nil
				}

				if err := cmd.Scan(&cart); err != nil {
					s.logger.Error("failed to scan session cart",
						zap.String("session_id", asessionID),
						zap.Any("value", cmd.Val()),
						zap.Error(err),
					)

					return nil
				}

				return nil
			})

			if err := g.Wait(); err != nil {
				return err //nolint:wrapcheck
			}

			acombine <- sessionStats{
				ID:        asessionID,
				TTL:       ttl,
				CreatedAt: cart.CreatedAt,
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
			g.Go(fetchStats(gCtx, combine, session))
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
			ID:        v.ID,
			TTL:       v.TTL,
			CreatedAt: v.CreatedAt,
		})
	}

	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return sessions, nil
}

func (s *Sessions) ListSessionPrivileges(ctx context.Context, sessionID string) ([]string, error) {
	const op = "Sessions.ListSessionPrivileges"

	ul, err := s.client.SMembers(ctx, s.keySession(sessionID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session privileges | %s:%w", op, err)
	}

	return ul, nil
}

func (s *Sessions) keyCart(sessionID string) string {
	return fmt.Sprintf("%s%s", spaceCarts, sessionID)
}

// Ключ сессии (с множеством привилегий)
func (s *Sessions) keySession(sessionID string) string {
	return fmt.Sprintf("%s%s", spaceSessions, sessionID)
}

// Ключ с сессиями пользователя
func (s *Sessions) keySessions(login string) string {
	return fmt.Sprintf("%s%s", space, login)
}
