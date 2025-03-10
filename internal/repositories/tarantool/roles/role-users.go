package tarantoolroles

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	tarantoolclient "github.com/vtievsky/auth-id/internal/repositories/tarantool/client"
)

func (s *Roles) GetRoleUsers(ctx context.Context, code string) ([]*models.RoleUser, error) {
	const op = "DbRoles.GetRoleUsers"

	role, err := s.GetRole(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	resp, err := s.c.Connection.Select(spaceRoleUser, "primary", 0, limit, tarantool.IterEq, tarantoolclient.Tuple{role.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get role users | %s:%w", op, err)
	}

	var roleUser tarantoolclient.RoleUser

	roleUsers := make([]*models.RoleUser, 0, len(resp.Tuples()))

	for _, tuple := range resp.Tuples() {
		roleUser = s.tupleToRoleUser(tuple)

		roleUsers = append(roleUsers, &models.RoleUser{
			RoleID:  roleUser.RoleID,
			UserID:  roleUser.UserID,
			DateIn:  roleUser.DateIn,
			DateOut: roleUser.DateOut,
		})
	}

	return roleUsers, nil
}

func (s *Roles) AddRoleUser(ctx context.Context, roleUser models.RoleUserCreated) error {
	const op = "DbRoles.AddRoleUser"

	roleUserCreated := tarantoolclient.RoleUserCreated{
		RoleID:  roleUser.RoleID,
		UserID:  roleUser.UserID,
		DateIn:  roleUser.DateIn,
		DateOut: roleUser.DateOut,
	}

	if _, err := s.c.Connection.Insert(spaceRoleUser, roleUserCreated.ToTuple()); err != nil {
		return fmt.Errorf("failed to add a role user | %s:%w", op, err)
	}

	return nil
}

func (s *Roles) UpdateRoleUser(ctx context.Context, roleUser models.RoleUserUpdated) error {
	const op = "DbRoles.UpdateRoleUser"

	roleUserUpdated := tarantoolclient.RoleUserUpdated{
		RoleID:  roleUser.RoleID,
		UserID:  roleUser.UserID,
		DateIn:  roleUser.DateIn,
		DateOut: roleUser.DateOut,
	}

	if _, err := s.c.Connection.Replace(spaceRoleUser, roleUserUpdated.ToTuple()); err != nil {
		return fmt.Errorf("failed to update a role user | %s:%w", op, err)
	}

	return nil
}

func (s *Roles) DeleteRoleUser(ctx context.Context, roleUser models.RoleUserDeleted) error {
	const op = "DbRoles.DeleteRoleUser"

	roleUserDeleted := tarantoolclient.RoleUserDeleted{
		RoleID: roleUser.RoleID,
		UserID: roleUser.UserID,
	}

	if _, err := s.c.Connection.Delete(spaceRoleUser, "pk", roleUserDeleted.ToTuple()); err != nil {
		return fmt.Errorf("failed to delete a role user | %s:%w", op, err)
	}

	return nil
}

func (s *Roles) tupleToRoleUser(tuple tarantoolclient.Tuple) tarantoolclient.RoleUser {
	return tarantoolclient.RoleUser{
		RoleID:  tuple[0].(uint64),                      //nolint:forcetypeassert
		UserID:  tuple[1].(uint64),                      //nolint:forcetypeassert
		DateIn:  time.Unix(int64(tuple[2].(uint64)), 0), //nolint:forcetypeassert,gosec
		DateOut: time.Unix(int64(tuple[3].(uint64)), 0), //nolint:forcetypeassert,gosec
	}
}
