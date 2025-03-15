package reposessions

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func (s *Sessions) Store(
	ctx context.Context,
	login, sessionID string,
	privileges []string,
	ttl time.Duration,
) error {
	const op = "Sessions.Store"

	keySession := s.keySession(sessionID)
	keySessions := s.keySessions(login)

	g, gCtx := errgroup.WithContext(ctx)

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
		if _, err := s.client.SAdd(gCtx, keySessions, keySession).Result(); err != nil {
			return fmt.Errorf("failed to add session to sessions list | %s:%w", op, err)
		}

		if _, err := s.client.Expire(gCtx, keySessions, ttl).Result(); err != nil {
			return fmt.Errorf("failed to set ttl sessions list | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}
