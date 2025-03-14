package redissessions

import (
	"context"
	"fmt"

	redisclient "github.com/vtievsky/auth-id/internal/repositories/redis/client"
)

const (
	limit = 5
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

// Ключ карты (с логином)
func (s *Sessions) keyCart(sessionID string) string {
	return fmt.Sprintf("crt:%s", sessionID)
}

func (s *Sessions) keyCarts(login string) string {
	return fmt.Sprintf("crt:%s", login)
}

// Ключ с набором привилегий
func (s *Sessions) keySession(sessionID string) string {
	return fmt.Sprintf("omo:%s", sessionID)
}

// Ключ с сессиями пользователя
func (s *Sessions) keySessions(login string) string {
	return fmt.Sprintf("omo:%s", login)
}

func (s *Sessions) fetchLoginSession(ctx context.Context, sessionID string) (string, error) {
	const op = "Sessions.fetchLoginSession"

	login, err := s.client.Get(ctx, s.keyCart(sessionID)).Result()
	if err != nil {
		return login, fmt.Errorf("failed to get session login | %s:%w", op, err)
	}

	return login, nil
}
