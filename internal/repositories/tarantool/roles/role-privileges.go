package tarantoolroles

import (
	"context"
	"fmt"

	"github.com/tarantool/go-tarantool"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	tarantoolclient "github.com/vtievsky/auth-id/internal/repositories/tarantool/client"
)

func (s *Roles) GetRolePrivileges(ctx context.Context, code string) ([]*models.RolePrivilege, error) {
	const op = "DbRoles.GetRolePrivileges"

	role, err := s.GetRole(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	resp, err := s.c.Connection.Select(spaceRolePrivilege, "pk", 0, limit, tarantool.IterEq, tarantoolclient.Tuple{role.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get role privileges | %s:%w", op, err)
	}

	roles := make([]*models.RolePrivilege, 0)

	for _, value := range resp.Tuples() {
		u := s.tupleToRolePrivilege(value)

		roles = append(roles, &models.RolePrivilege{
			RoleID:      int(u.RoleID),      //nolint:gosec
			PrivilegeID: int(u.PrivilegeID), //nolint:gosec
			Allowed:     u.Allowed,
		})
	}

	return roles, nil
}

func (s *Roles) AddRolePrivilege(ctx context.Context, rolePrivilege models.RolePrivilegeCreated) error {
	const op = "DbRoles.AddRolePrivilege"

	if _, err := s.c.Connection.Insert(
		spaceRolePrivilege,
		s.rolePrivilegeCreatedToTuple(
			tarantoolclient.RolePrivilegeCreated{
				RoleID:      uint64(rolePrivilege.RoleID),      //nolint:gosec
				PrivilegeID: uint64(rolePrivilege.PrivilegeID), //nolint:gosec
				Allowed:     rolePrivilege.Allowed,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to add a role privilege | %s:%w", op, err)
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

func (s *Roles) rolePrivilegeCreatedToTuple(data tarantoolclient.RolePrivilegeCreated) tarantoolclient.Tuple {
	return tarantoolclient.Tuple{
		data.RoleID,
		data.PrivilegeID,
		data.Allowed,
	}
}
