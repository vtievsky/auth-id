package redissessions

import (
	"context"
	"fmt"
	"time"

	redisclient "github.com/vtievsky/auth-id/internal/repositories/redis/client"
	"golang.org/x/sync/errgroup"
)

type SessionsOpts struct {
	Client *redisclient.Client
}

type Sessions struct {
	client *redisclient.Client
}

func New(opts *SessionsOpts) *Sessions {
	return &Sessions{
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

func (s *Sessions) Store(
	ctx context.Context,
	login, sessionID string,
	privileges []string,
	ttl time.Duration,
) error {
	const op = "Sessions.Store"

	keyCart := s.keyCart(sessionID)
	keySession := s.keySession(sessionID)
	keyLoginSessions := s.keyLoginSessions(login)

	g, gCtx := errgroup.WithContext(ctx)

	// Сохранение карты пользователя
	g.Go(func() error {
		if _, err := s.client.Set(gCtx, keyCart, login, ttl).Result(); err != nil {
			return fmt.Errorf("failed to add session cart | %s:%w", op, err)
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
		if _, err := s.client.SAdd(gCtx, keyLoginSessions, keySession).Result(); err != nil {
			return fmt.Errorf("failed to add session to sessions list | %s:%w", op, err)
		}

		if _, err := s.client.Expire(gCtx, keyLoginSessions, ttl).Result(); err != nil {
			return fmt.Errorf("failed to set ttl sessions list | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

func (s *Sessions) Delete(ctx context.Context, sessionID string) error {
	const op = "Sessions.Delete"

	keyCart := s.keyCart(sessionID)

	login, err := s.client.Get(ctx, keyCart).Result()
	if err != nil {
		return fmt.Errorf("failed to get session cart | %s:%w", op, err)
	}

	keySession := s.keySession(sessionID)
	keyLoginSessions := s.keyLoginSessions(login)

	g, gCtx := errgroup.WithContext(ctx)

	// Удаление карты пользователя
	g.Go(func() error {
		if _, err := s.client.Del(gCtx, keyCart).Result(); err != nil {
			return fmt.Errorf("failed to remove session cart | %s:%w", op, err)
		}

		return nil
	})

	// Удаление списка привилегий сессии
	g.Go(func() error {
		if _, err := s.client.Del(gCtx, keySession).Result(); err != nil {
			return fmt.Errorf("failed to remove session privileges | %s:%w", op, err)
		}

		return nil
	})

	// Удаление сессии из списка сессий пользователя
	g.Go(func() error {
		if _, err := s.client.SRem(gCtx, keyLoginSessions, keySession).Result(); err != nil {
			return fmt.Errorf("failed to remove session from sessions list | %s:%w", op, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

// Ключ карты (с логином)
func (s *Sessions) keyCart(login string) string {
	return fmt.Sprintf("pok:%s", login)
}

// Ключ с набором привилегий
func (s *Sessions) keySession(sessionID string) string {
	return fmt.Sprintf("omo:%s", sessionID)
}

// Ключ с сессиями пользователя
func (s *Sessions) keyLoginSessions(login string) string {
	return fmt.Sprintf("omo:%s", login)
}
