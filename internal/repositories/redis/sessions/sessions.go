package redissessions

import (
	"context"
	"fmt"
	"time"

	redisclient "github.com/vtievsky/auth-id/internal/repositories/redis/client"
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

func (s *Sessions) Store(
	ctx context.Context,
	login, sessionID string,
	privileges []string,
	ttl time.Duration,
) error {
	const op = "Sessions.Store"

	// Сохранение списка привилегий сессии
	keySession := s.keySession(sessionID)

	if _, err := s.client.SAdd(ctx, keySession, privileges).Result(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if _, err := s.client.Expire(ctx, keySession, ttl).Result(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	// Добавление сессии в список сессий пользователя
	keySessions := s.keySessions(login)

	if _, err := s.client.SAdd(ctx, keySessions, keySession).Result(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if _, err := s.client.Expire(ctx, keySessions, ttl).Result(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Sessions) keySession(sessionID string) string {
	return fmt.Sprintf("%s:%s", space, sessionID)
}

func (s *Sessions) keySessions(login string) string {
	return fmt.Sprintf("%s:%s", space, login)
}
