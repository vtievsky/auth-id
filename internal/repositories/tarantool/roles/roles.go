package tarantoolroles

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool"
	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	tarantoolclient "github.com/vtievsky/auth-id/internal/repositories/tarantool/client"
)

const (
	spaceRole          = "role"
	spaceRolePrivilege = "role_privilege"
	spaceRoleUser      = "role_user"
	limit              = 25
)

type RolesOpts struct {
	Client *tarantoolclient.Client
}

type Roles struct {
	c *tarantoolclient.Client
}

func New(opts *RolesOpts) *Roles {
	return &Roles{
		c: opts.Client,
	}
}

func (s *Roles) GetRole(ctx context.Context, code string) (*models.Role, error) {
	const op = "DbRoles.GetRole"

	resp, err := s.c.Connection.Select(spaceRole, "secondary", 0, 1, tarantool.IterEq, tarantoolclient.Tuple{code})
	if err != nil {
		return nil, fmt.Errorf("failed to get role | %s:%w", op, err)
	}

	if len(resp.Tuples()) < 1 {
		return nil, fmt.Errorf("failed to get role | %s:%w", op, dberrors.ErrRoleNotFound)
	}

	role := s.tupleToRole(resp.Tuples()[0])

	return &models.Role{
		ID:          int(role.ID), //nolint:gosec
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		Blocked:     role.Blocked,
	}, nil
}

func (s *Roles) GetRoles(ctx context.Context) ([]*models.Role, error) {
	const op = "DbRoles.GetRoles"

	resp, err := s.c.Connection.Select(spaceRole, "pk", 0, limit, tarantool.IterAll, tarantoolclient.Tuple{})
	if err != nil {
		return nil, fmt.Errorf("failed to get roles | %s:%w", op, err)
	}

	var role tarantoolclient.Role

	roles := make([]*models.Role, 0, len(resp.Tuples()))

	for _, tuple := range resp.Tuples() {
		role = s.tupleToRole(tuple)

		roles = append(roles, &models.Role{
			ID:          int(role.ID), //nolint:gosec
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			Blocked:     role.Blocked,
		})
	}

	return roles, nil
}

func (s *Roles) CreateRole(ctx context.Context, role models.RoleCreated) (*models.Role, error) {
	const op = "DbRoles.CreateRole"

	roleCreated := tarantoolclient.RoleCreated{
		Code:        fmt.Sprintf("r%d", time.Now().Unix()),
		Name:        role.Name,
		Description: role.Description,
		Blocked:     role.Blocked,
	}

	if _, err := s.c.Connection.Insert(spaceRole, roleCreated.ToTuple()); err != nil {
		return nil, fmt.Errorf("failed to create role | %s:%w", op, err)
	}

	return s.GetRole(ctx, roleCreated.Code)
}

func (s *Roles) UpdateRole(ctx context.Context, role models.RoleUpdated) (*models.Role, error) {
	const op = "DbRoles.UpdateRole"

	u, err := s.GetRole(ctx, role.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to update role | %s:%w", op, err)
	}

	roleUpdated := tarantoolclient.RoleUpdated{
		ID:          uint64(u.ID), //nolint:gosec
		Name:        role.Name,
		Description: u.Description,
		Blocked:     role.Blocked,
	}

	if _, err := s.c.Connection.Replace(spaceRole, roleUpdated.ToTuple()); err != nil {
		return nil, fmt.Errorf("failed to update user | %s:%w", op, err)
	}

	return s.GetRole(ctx, role.Code)
}

func (s *Roles) DeleteRole(ctx context.Context, code string) error {
	const op = "DbRoles.DeleteRole"

	if _, err := s.GetRole(ctx, code); err != nil {
		return fmt.Errorf("failed to delete role | %s:%w", op, err)
	}

	if _, err := s.c.Connection.Delete(spaceRole, "secondary", tarantoolclient.Tuple{code}); err != nil {
		return fmt.Errorf("failed to delete role | %s:%w", op, err)
	}

	return nil
}

func (s *Roles) tupleToRole(tuple tarantoolclient.Tuple) tarantoolclient.Role {
	return tarantoolclient.Role{
		ID:          tuple[0].(uint64), //nolint:forcetypeassert
		Code:        tuple[1].(string), //nolint:forcetypeassert
		Name:        tuple[2].(string), //nolint:forcetypeassert
		Description: tuple[3].(string), //nolint:forcetypeassert
		Blocked:     tuple[4].(bool),   //nolint:forcetypeassert
	}
}
