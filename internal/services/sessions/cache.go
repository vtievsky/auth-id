package sessionsvc

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"go.uber.org/zap"
)

func (s *SessionSvc) Search(ctx context.Context, sessionID, privilegeCode string) error {
	const op = "SessionSvc.Search"

	privileges, err := s.cacheByID.Get(ctx, sessionID,
		func(ctx context.Context) (map[string][]string, error) {
			privileges, err := s.storage.ListSessionPrivileges(ctx, sessionID)

			switch {
			case errors.Is(err, nil):
				s.logger.Debug("session privileges has been synchronized",
					zap.String("session_id", sessionID),
					zap.Int("num", len(privileges)),
				)

				return map[string][]string{
					sessionID: privileges,
				}, nil
			default:
				return nil, err //nolint:wrapcheck
			}
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

	s.logger.Error("failed to search session privilege",
		zap.String("session_id", sessionID),
		zap.String("privilege_code", privilegeCode),
		zap.Error(ErrSessionPrivilegeNotFound),
	)

	return ErrSessionPrivilegeNotFound
}
