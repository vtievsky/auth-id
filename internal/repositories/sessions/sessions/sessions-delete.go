package redissessions

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// Удаляет указанную сессию пользователя
func (s *Sessions) Delete(ctx context.Context, login, sessionID string) error {
	const op = "Sessions.Delete"

	keySession := s.keySession(sessionID)
	keySessions := s.keySessions(login)

	g, gCtx := errgroup.WithContext(ctx)

	// Удаление списка привилегий сессии
	g.Go(func() error {
		if _, err := s.client.Del(gCtx, keySession).Result(); err != nil {
			return fmt.Errorf("failed to delete session privileges | %s:%w", op, err)
		}

		return nil
	})

	// Удаление сессии из списка сессий пользователя
	g.Go(func() error {
		if _, err := s.client.SRem(gCtx, keySessions, keySession).Result(); err != nil {
			return fmt.Errorf("failed to delete session from sessions list | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}
