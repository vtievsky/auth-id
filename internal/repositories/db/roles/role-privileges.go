package tarantoolroles

import (
	"context"
	"fmt"

	"github.com/tarantool/go-tarantool"
	tarantoolclient "github.com/vtievsky/auth-id/internal/repositories/db/client"
	"github.com/vtievsky/auth-id/internal/repositories/models"
)

func (s *Roles) GetRolePrivileges(ctx context.Context, code string) ([]*models.RolePrivilege, error) {
	const op = "DbRoles.GetRolePrivileges"

	role, err := s.GetRole(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	resp, err := s.c.Connection.Select(spaceRolePrivilege, "primary", 0, limit, tarantool.IterEq, tarantoolclient.Tuple{role.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get role privileges | %s:%w", op, err)
	}

	var rolePrivilege tarantoolclient.RolePrivilege

	rolePrivileges := make([]*models.RolePrivilege, 0, len(resp.Tuples()))

	for _, tuple := range resp.Tuples() {
		rolePrivilege = s.tupleToRolePrivilege(tuple)

		rolePrivileges = append(rolePrivileges, &models.RolePrivilege{
			RoleID:      rolePrivilege.RoleID,
			PrivilegeID: rolePrivilege.PrivilegeID,
			Allowed:     rolePrivilege.Allowed,
		})
	}

	return rolePrivileges, nil
}

func (s *Roles) AddRolePrivilege(ctx context.Context, rolePrivilege models.RolePrivilegeCreated) error {
	const op = "DbRoles.AddRolePrivilege"

	rolePrivilegeCreated := tarantoolclient.RolePrivilegeCreated{
		RoleID:      rolePrivilege.RoleID,
		PrivilegeID: rolePrivilege.PrivilegeID,
		Allowed:     rolePrivilege.Allowed,
	}

	if _, err := s.c.Connection.Insert(spaceRolePrivilege, rolePrivilegeCreated.ToTuple()); err != nil {
		return fmt.Errorf("failed to add a role privilege | %s:%w", op, err)
	}

	return nil
}

func (s *Roles) UpdateRolePrivilege(ctx context.Context, rolePrivilege models.RolePrivilegeUpdated) error {
	const op = "DbRoles.UpdateRolePrivilege"

	rolePrivilegeUpdated := tarantoolclient.RolePrivilegeUpdated{
		RoleID:      rolePrivilege.RoleID,
		PrivilegeID: rolePrivilege.PrivilegeID,
		Allowed:     rolePrivilege.Allowed,
	}

	if _, err := s.c.Connection.Replace(spaceRolePrivilege, rolePrivilegeUpdated.ToTuple()); err != nil {
		return fmt.Errorf("failed to update a role privilege | %s:%w", op, err)
	}

	return nil
}

func (s *Roles) DeleteRolePrivilege(ctx context.Context, rolePrivilege models.RolePrivilegeDeleted) error {
	const op = "DbRoles.DeleteRolePrivilege"

	rolePrivilegeDeleted := tarantoolclient.RolePrivilegeDeleted{
		RoleID:      rolePrivilege.RoleID,
		PrivilegeID: rolePrivilege.PrivilegeID,
	}

	if _, err := s.c.Connection.Delete(spaceRolePrivilege, "pk", rolePrivilegeDeleted.ToTuple()); err != nil {
		return fmt.Errorf("failed to delete a role privilege | %s:%w", op, err)
	}

	return nil
}

func (s *Roles) tupleToRolePrivilege(tuple tarantoolclient.Tuple) tarantoolclient.RolePrivilege {
	return tarantoolclient.RolePrivilege{
		RoleID:      tuple[0].(uint64), //nolint:forcetypeassert
		PrivilegeID: tuple[1].(uint64), //nolint:forcetypeassert
		Allowed:     tuple[2].(bool),   //nolint:forcetypeassert
	}
}
