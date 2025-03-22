package tarantoolusers

import (
	"context"
	"fmt"

	"github.com/tarantool/go-tarantool"
	clienttarantool "github.com/vtievsky/auth-id/internal/repositories/db/client/tarantool"
	"github.com/vtievsky/auth-id/internal/repositories/models"
)

func (s *Users) GetUserRoles(ctx context.Context, login string, pageSize, offset uint32) ([]*models.UserRole, error) {
	const op = "DbUsers.GetUserRoles"

	user, err := s.GetUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	resp, err := s.c.Connection.Select(spaceUserRole, "secondary", offset, pageSize, tarantool.IterEq, clienttarantool.Tuple{user.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles | %s:%w", op, err)
	}

	var userRole clienttarantool.UserRole

	userRoles := make([]*models.UserRole, 0, len(resp.Tuples()))

	for _, tuple := range resp.Tuples() {
		userRole = clienttarantool.Tuple(tuple).ToUserRole()

		userRoles = append(userRoles, &models.UserRole{
			RoleID:  userRole.RoleID,
			UserID:  userRole.UserID,
			DateIn:  userRole.DateIn,
			DateOut: userRole.DateOut,
		})
	}

	return userRoles, nil
}
