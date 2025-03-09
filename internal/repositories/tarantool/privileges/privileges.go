package tarantoolprivileges

import (
	"context"
	"fmt"

	"github.com/tarantool/go-tarantool"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	tarantoolclient "github.com/vtievsky/auth-id/internal/repositories/tarantool/client"
)

const (
	space = "privilege"
	limit = 25
)

type PrivilegesOpts struct {
	Client *tarantoolclient.Client
}

type Privileges struct {
	c *tarantoolclient.Client
}

func New(opts *PrivilegesOpts) *Privileges {
	return &Privileges{
		c: opts.Client,
	}
}

func (s *Privileges) GetPrivileges(ctx context.Context) ([]*models.Privilege, error) {
	const op = "DbPrivileges.GetPrivileges"

	resp, err := s.c.Connection.Select(space, "pk", 0, limit, tarantool.IterAll, tarantoolclient.Tuple{})
	if err != nil {
		return nil, fmt.Errorf("failed to get privileges | %s:%w", op, err)
	}

	privileges := make([]*models.Privilege, 0)

	for _, value := range resp.Tuples() {
		u := s.tupleToPrivilege(value)

		privileges = append(privileges, &models.Privilege{
			ID:          int(u.ID), //nolint:gosec
			Code:        u.Code,
			Name:        u.Name,
			Description: u.Description,
		})
	}

	return privileges, nil
}

func (s *Privileges) tupleToPrivilege(tuple tarantoolclient.Tuple) tarantoolclient.Privilege {
	return tarantoolclient.Privilege{
		ID:          tuple[0].(uint64), //nolint:forcetypeassert
		Code:        tuple[1].(string), //nolint:forcetypeassert
		Name:        tuple[2].(string), //nolint:forcetypeassert
		Description: tuple[3].(string), //nolint:forcetypeassert
	}
}
