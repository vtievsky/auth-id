package tarantoolusers

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool"
	clienttarantool "github.com/vtievsky/auth-id/internal/repositories/db/client/tarantool"
	"github.com/vtievsky/auth-id/internal/repositories/models"
)

func (s *Users) GetUserRoles(ctx context.Context, login string) ([]*models.UserRole, error) {
	const op = "DbUsers.GetUserRoles"

	user, err := s.GetUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	resp, err := s.c.Connection.Select(spaceUserRole, "secondary", 0, limit, tarantool.IterEq, clienttarantool.Tuple{user.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles | %s:%w", op, err)
	}

	var roleUser clienttarantool.UserRole

	roleUsers := make([]*models.UserRole, 0, len(resp.Tuples()))

	for _, tuple := range resp.Tuples() {
		roleUser = s.tupleToUserRole(tuple)

		roleUsers = append(roleUsers, &models.UserRole{
			RoleID:  roleUser.RoleID,
			UserID:  roleUser.UserID,
			DateIn:  roleUser.DateIn,
			DateOut: roleUser.DateOut,
		})
	}

	return roleUsers, nil
}

func (s *Users) tupleToUserRole(tuple clienttarantool.Tuple) clienttarantool.UserRole {
	return clienttarantool.UserRole{
		RoleID:  tuple[0].(uint64),                      //nolint:forcetypeassert
		UserID:  tuple[1].(uint64),                      //nolint:forcetypeassert
		DateIn:  time.Unix(int64(tuple[2].(uint64)), 0), //nolint:forcetypeassert,gosec
		DateOut: time.Unix(int64(tuple[3].(uint64)), 0), //nolint:forcetypeassert,gosec
	}
}
