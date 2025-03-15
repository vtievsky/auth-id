package redissessions

import (
	"context"
)

// Удаляет указанную сессию пользователя
func (s *Sessions) Delete(ctx context.Context, sessionID string) error {
	const op = "Sessions.Delete"

	// login, err := s.fetchLoginSession(ctx, sessionID)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete session | %s:%w", op, err)
	// }

	// keyCart := s.keyCart(sessionID)
	// keyCarts := s.keyCarts(login)
	// keySession := s.keySession(sessionID)
	// keySessions := s.keySessions(login)

	// g, gCtx := errgroup.WithContext(ctx)

	// // Удаление карты сессии
	// g.Go(func() error {
	// 	if _, err := s.client.Del(gCtx, keyCart).Result(); err != nil {
	// 		return fmt.Errorf("failed to delete session cart | %s:%w", op, err)
	// 	}

	// 	return nil
	// })

	// // Удаление карты из списка карт сессии
	// g.Go(func() error {
	// 	if _, err := s.client.SRem(gCtx, keyCarts, keyCart).Result(); err != nil {
	// 		return fmt.Errorf("failed to delete session cart from carts list | %s:%w", op, err)
	// 	}

	// 	return nil
	// })

	// // Удаление списка привилегий сессии
	// g.Go(func() error {
	// 	if _, err := s.client.Del(gCtx, keySession).Result(); err != nil {
	// 		return fmt.Errorf("failed to delete session privileges | %s:%w", op, err)
	// 	}

	// 	return nil
	// })

	// // Удаление сессии из списка сессий пользователя
	// g.Go(func() error {
	// 	if _, err := s.client.SRem(gCtx, keySessions, keySession).Result(); err != nil {
	// 		return fmt.Errorf("failed to delete session from sessions list | %s:%w", op, err)
	// 	}

	// 	return nil
	// })

	// if err := g.Wait(); err != nil {
	// 	return err //nolint:wrapcheck
	// }

	return nil
}
