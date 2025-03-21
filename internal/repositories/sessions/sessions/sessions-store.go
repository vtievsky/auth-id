package reposessions

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

type SessionCart struct {
	ID        string    `redis:"id"`
	CreatedAt time.Time `redis:"created_at"`
}

func (s *Sessions) Store(
	ctx context.Context,
	login, sessionID string,
	privileges []string,
	ttl time.Duration,
) error {
	const op = "Sessions.Store"

	if len(privileges) < 1 {
		return ErrSessionPrivilegesEmpty
	}

	keyCart := s.keyCart(sessionID)
	keySession := s.keySession(sessionID)
	keySessions := s.keySessions(login)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if _, err := s.client.HMSet(ctx, keyCart, SessionCart{
			ID:        sessionID,
			CreatedAt: time.Now(),
		}).Result(); err != nil {
			return fmt.Errorf("failed to add session cart | %s:%w", op, err)
		}

		if _, err := s.client.Expire(gCtx, keyCart, ttl).Result(); err != nil {
			return fmt.Errorf("failed to set ttl session cart | %s:%w", op, err)
		}

		return nil
	})

	// Сохранение списка привилегий сессии
	g.Go(func() error {
		if _, err := s.client.SAdd(gCtx, keySession, privileges).Result(); err != nil {
			return fmt.Errorf("failed to add session privileges | %s:%w", op, err)
		}

		if _, err := s.client.Expire(gCtx, keySession, ttl).Result(); err != nil {
			return fmt.Errorf("failed to set ttl session privileges | %s:%w", op, err)
		}

		return nil
	})

	// Добавление сессии в список сессий пользователя
	g.Go(func() error {
		if _, err := s.client.SAdd(gCtx, keySessions, sessionID).Result(); err != nil {
			return fmt.Errorf("failed to add session to sessions list | %s:%w", op, err)
		}

		if _, err := s.client.Expire(gCtx, sessionID, ttl).Result(); err != nil {
			return fmt.Errorf("failed to set ttl sessions list | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}
