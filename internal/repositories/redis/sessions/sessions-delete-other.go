package redissessions

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/sync/errgroup"
)

// Удаляет все сессии пользователя, кроме указанной
func (s *Sessions) DeleteOther(ctx context.Context, sessionID string) error {
	const op = "Sessions.DeleteOther"

	login, err := s.fetchLoginSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete other sessions | %s:%w", op, err)
	}

	keyCart := s.keyCart(sessionID)
	keyCarts := s.keyCarts(login)
	keySession := s.keySession(sessionID)
	keySessions := s.keySessions(login)

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(limit)

	// Извлекаем карты пользователя
	carts, err := s.client.SMembers(ctx, keyCarts).Result()
	if err != nil {
		return fmt.Errorf("failed to get carts list | %s:%w", op, err)
	}

	for _, cart := range carts {
		if strings.EqualFold(keyCart, cart) {
			continue
		}

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
		if strings.EqualFold(keySession, session) {
			continue
		}

		// Удаление карты сессии
		g.Go(func() error {
			if _, err := s.client.Del(gCtx, session).Result(); err != nil {
				return fmt.Errorf("failed to delete session | %s:%w", op, err)
			}

			return nil
		})

		// Удаление из списка карт сессии
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
