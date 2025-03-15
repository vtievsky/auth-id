package tarantoolprivileges

import (
	"context"
	"fmt"

	"github.com/tarantool/go-tarantool"
	clienttarantool "github.com/vtievsky/auth-id/internal/repositories/db/client/tarantool"
	"github.com/vtievsky/auth-id/internal/repositories/models"
)

const (
	space = "privilege"
	limit = 25
)

type PrivilegesOpts struct {
	Client *clienttarantool.Client
}

type Privileges struct {
	c *clienttarantool.Client
}

func New(opts *PrivilegesOpts) *Privileges {
	return &Privileges{
		c: opts.Client,
	}
}

func (s *Privileges) GetPrivileges(ctx context.Context) ([]*models.Privilege, error) {
	const op = "DbPrivileges.GetPrivileges"

	resp, err := s.c.Connection.Select(space, "pk", 0, limit, tarantool.IterAll, clienttarantool.Tuple{})
	if err != nil {
		return nil, fmt.Errorf("failed to get privileges | %s:%w", op, err)
	}

	var value clienttarantool.Privilege

	privileges := make([]*models.Privilege, 0, len(resp.Tuples()))

	for _, tuple := range resp.Tuples() {
		value = s.tupleToPrivilege(tuple)

		privileges = append(privileges, &models.Privilege{
			ID:          value.ID,
			Code:        value.Code,
			Name:        value.Name,
			Description: value.Description,
		})
	}

	return privileges, nil
}

func (s *Privileges) tupleToPrivilege(tuple clienttarantool.Tuple) clienttarantool.Privilege {
	return clienttarantool.Privilege{
		ID:          tuple[0].(uint64), //nolint:forcetypeassert
		Code:        tuple[1].(string), //nolint:forcetypeassert
		Name:        tuple[2].(string), //nolint:forcetypeassert
		Description: tuple[3].(string), //nolint:forcetypeassert
	}
}
