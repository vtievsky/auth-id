package redissessions

import (
	"context"
	"fmt"
	"time"

	redisclient "github.com/vtievsky/auth-id/internal/repositories/redis/client"
	"golang.org/x/sync/errgroup"
)

const (
	space = "ses"
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

	keySession := s.keySession(sessionID)
	keyLoginSessions := s.keyLoginSessions(login)

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

func (s *Sessions) Delete(ctx context.Context, login, sessionID string) error {
	const op = "Sessions.Delete"

	keySession := s.keySession(sessionID)
	keyLoginSessions := s.keyLoginSessions(login)

	g, gCtx := errgroup.WithContext(ctx)

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

func (s *Sessions) keySession(sessionID string) string {
	return fmt.Sprintf("%s:%s", space, sessionID)
}

func (s *Sessions) keyLoginSessions(login string) string {
	return fmt.Sprintf("%s:%s", space, login)
}
