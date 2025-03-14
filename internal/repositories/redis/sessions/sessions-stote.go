package redissessions

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

	keyCart := s.keyCart(sessionID)
	keyCarts := s.keyCarts(login)
	keySession := s.keySession(sessionID)
	keySessions := s.keySessions(login)

	g, gCtx := errgroup.WithContext(ctx)

	// Сохранение карты пользователя
	g.Go(func() error {
		if _, err := s.client.Set(gCtx, keyCart, login, ttl).Result(); err != nil {
			return fmt.Errorf("failed to add session cart | %s:%w", op, err)
		}

		return nil
	})

	// Сохранение списка карт сессии
	g.Go(func() error {
		if _, err := s.client.SAdd(gCtx, keyCarts, keyCart).Result(); err != nil {
			return fmt.Errorf("failed to add cart to carts list| %s:%w", op, err)
		}

		if _, err := s.client.Expire(gCtx, keyCarts, ttl).Result(); err != nil {
			return fmt.Errorf("failed to set ttl carts list | %s:%w", op, err)
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
