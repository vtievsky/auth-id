package reposessions

import (
	"context"
	"fmt"
	"slices"
	"time"

	clientredis "github.com/vtievsky/auth-id/internal/repositories/sessions/client/redis"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	space         = "uev:"
	spaceCarts    = "pak:"
	spaceSessions = "omo:"
	threadsLimit  = 5
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

func (s *Sessions) Get(ctx context.Context, sessionID string) (*SessionCart, error) {
	const op = "Sessions.Get"

	cmd := s.client.HGetAll(ctx, s.keyCart(sessionID))

	switch {
	case cmd.Err() != nil:
		s.logger.Error("failed to get session cart",
			zap.String("session_id", sessionID),
			zap.Error(cmd.Err()),
		)

		return nil, fmt.Errorf("failed to get session cart | %s:%w", op, cmd.Err())
	case len(cmd.Val()) == 0:
		s.logger.Error("failed to get session cart",
			zap.String("session_id", sessionID),
			zap.Error(ErrSessionCartNotFound),
		)

		return nil, fmt.Errorf("failed to get session cart | %s:%w", op, ErrSessionCartNotFound)
	}

	var cart SessionCart

	if err := cmd.Scan(&cart); err != nil {
		s.logger.Error("failed to scan session cart",
			zap.String("session_id", sessionID),
			zap.Any("value", cmd.Val()),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get session cart | %s:%w", op, err)
	}

	return &cart, nil
}

func (s *Sessions) List(ctx context.Context, login string, pageSize, offset uint32) ([]*Session, error) {
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

			if ttl < 0 {
				return nil
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
	g.SetLimit(threadsLimit + 1)

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

	// Реализация пагинации
	_len := len(sessions)
	_offset := int(offset)
	_pageSize := int(pageSize)
	_bucketSize := _offset + _pageSize

	switch {
	case _len < _offset:
		return []*Session{}, nil
	case _len < _bucketSize:
		_bucketSize = _len
	}

	// "Постраничный" вывод
	slices.SortStableFunc(sessions, func(a, b *Session) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return -1
		}

		return 1
	})

	sessions = sessions[_offset:_bucketSize]

	return sessions, nil
}

func (s *Sessions) ListSessionPrivileges(ctx context.Context, sessionID string, pageSize, offset uint32) ([]string, error) {
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
