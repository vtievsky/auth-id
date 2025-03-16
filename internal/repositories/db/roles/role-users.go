package tarantoolroles

import (
	"context"
	"fmt"

	"github.com/tarantool/go-tarantool"
	clienttarantool "github.com/vtievsky/auth-id/internal/repositories/db/client/tarantool"
	"github.com/vtievsky/auth-id/internal/repositories/models"
)

func (s *Roles) GetRoleUsers(ctx context.Context, code string) ([]*models.RoleUser, error) {
	const op = "DbRoles.GetRoleUsers"

	role, err := s.GetRole(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	resp, err := s.c.Connection.Select(spaceRoleUser, "primary", 0, limit, tarantool.IterEq, clienttarantool.Tuple{role.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get role users | %s:%w", op, err)
	}

	var roleUser clienttarantool.RoleUser

	roleUsers := make([]*models.RoleUser, 0, len(resp.Tuples()))

	for _, tuple := range resp.Tuples() {
		roleUser = clienttarantool.Tuple(tuple).ToRoleUser()

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

	roleUserCreated := clienttarantool.RoleUserCreated{
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

	roleUserUpdated := clienttarantool.RoleUserUpdated{
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

	roleUserDeleted := clienttarantool.RoleUserDeleted{
		RoleID: roleUser.RoleID,
		UserID: roleUser.UserID,
	}

	if _, err := s.c.Connection.Delete(spaceRoleUser, "pk", roleUserDeleted.ToTuple()); err != nil {
		return fmt.Errorf("failed to delete a role user | %s:%w", op, err)
	}

	return nil
}
