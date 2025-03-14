package redissessions

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// Удаляет все сессии того пользователя, чья сессия указана
func (s *Sessions) DeleteAll(ctx context.Context, sessionID string) error {
	const op = "Sessions.DeleteAll"

	login, err := s.fetchLoginSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete all sessions | %s:%w", op, err)
	}

	keyCarts := s.keyCarts(login)
	keySessions := s.keySessions(login)

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(limit)

	// Извлекаем карты пользователя
	carts, err := s.client.SMembers(ctx, keyCarts).Result()
	if err != nil {
		return fmt.Errorf("failed to get carts list | %s:%w", op, err)
	}

	for _, cart := range carts {
		// Удаление карты сессии
		g.Go(func() error {
			if _, err := s.client.Del(gCtx, cart).Result(); err != nil {
				return fmt.Errorf("failed to delete session cart | %s:%w", op, err)
			}

			return nil
		})

		// Удаление из списка карт сессии
		g.Go(func() error {
			if _, err := s.client.SRem(gCtx, keyCarts, cart).Result(); err != nil {
				return fmt.Errorf("failed to delete session cart from carts list | %s:%w", op, err)
			}

			return nil
		})
	}

	// Извлекаем сессии пользователя
	sessions, err := s.client.SMembers(ctx, keySessions).Result()
	if err != nil {
		return fmt.Errorf("failed to get sessions list | %s:%w", op, err)
	}

	for _, session := range sessions {
		// Удаление сессии
		g.Go(func() error {
			if _, err := s.client.Del(gCtx, session).Result(); err != nil {
				return fmt.Errorf("failed to delete session | %s:%w", op, err)
			}

			return nil
		})

		// Удаление из списка сессий
		g.Go(func() error {
			if _, err := s.client.SRem(gCtx, keySessions, session).Result(); err != nil {
				return fmt.Errorf("failed to delete session from sessions list | %s:%w", op, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}
