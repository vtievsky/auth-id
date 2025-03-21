package sessionsvc

import (
	"context"
	"errors"
	"fmt"
	"slices"

	reposessions "github.com/vtievsky/auth-id/internal/repositories/sessions/sessions"
	"go.uber.org/zap"
)

func (s *SessionSvc) Find(ctx context.Context, sessionID, privilegeCode string) error {
	const op = "SessionSvc.Find"

	privileges, err := s.cacheByID.Get(ctx, sessionID,
		func(ctx context.Context) (map[string][]string, error) {
			err := s.storage.Find(ctx, sessionID, privilegeCode)

			var resp map[string][]string

			switch {
			case errors.Is(err, nil):
				resp = map[string][]string{
					sessionID: {privilegeCode},
				}
			case errors.Is(err, reposessions.ErrSessionPrivilegeNotFound):
				resp = map[string][]string{
					sessionID: {},
				}
			default:
				return nil, err //nolint:wrapcheck
			}

			s.logger.Debug("privileges has been synchronized",
				zap.String("session_id", sessionID),
				zap.String("privilege_code", privilegeCode),
			)

			return resp, nil
		})
	if err != nil {
		s.logger.Error("failed to search session privilege",
			zap.String("session_id", sessionID),
			zap.String("privilege_code", privilegeCode),
			zap.Error(err),
		)

		return fmt.Errorf("failed to search session privilege | %s:%w", op, err)
	}

	if slices.Contains(privileges, privilegeCode) {
		return nil
	}

	return ErrSessionPrivilegeNotFound
}
