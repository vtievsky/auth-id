package tarantoolusers

import (
	"context"
	"fmt"

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
		roleUser = clienttarantool.Tuple(tuple).ToUserRole()

		roleUsers = append(roleUsers, &models.UserRole{
			RoleID:  roleUser.RoleID,
			UserID:  roleUser.UserID,
			DateIn:  roleUser.DateIn,
			DateOut: roleUser.DateOut,
		})
	}

	return roleUsers, nil
}
