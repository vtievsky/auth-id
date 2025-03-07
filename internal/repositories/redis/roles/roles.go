package redisroles

import (
	"context"
	"fmt"
	"strings"
	"time"

	dberrors "github.com/vtievsky/auth-id/internal/repositories"
	"github.com/vtievsky/auth-id/internal/repositories/models"
	redisclient "github.com/vtievsky/auth-id/internal/repositories/redis/client"
)

const (
	space = "rol"
)

type RolesOpts struct {
	Client *redisclient.Client
}

type Roles struct {
	client *redisclient.Client
}

func New(opts *RolesOpts) *Roles {
	return &Roles{
		client: opts.Client,
	}
}

func (s *Roles) GetRole(ctx context.Context, code string) (*models.Role, error) {
	const op = "DbRoles.GetRole"

	cmd := s.client.HGetAll(ctx, s.codeToKey(code))

	switch {
	case cmd.Err() != nil:
		return nil, fmt.Errorf("failed to get role | %s:%w", op, cmd.Err())
	case len(cmd.Val()) < 1:
		return nil, fmt.Errorf("failed to get role | %s:%w", op, dberrors.ErrRoleNotFound)
	}

	var value redisclient.Role

	err := cmd.Scan(&value)
	if err != nil {
		return nil, fmt.Errorf("failed to get role | %s:%w", op, dberrors.ErrRoleScan)
	}

	return &models.Role{
		ID:          value.ID,
		Code:        value.Code,
		Name:        value.Name,
		Description: value.Description,
		Blocked:     value.Blocked,
	}, nil
}

func (s *Roles) GetRoles(ctx context.Context) ([]*models.Role, error) {
	const op = "DbRoles.GetRoles"

	ul, err := s.client.Keys(ctx, s.space()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get roles | %s:%w", op, dberrors.ErrRoleScan)
	}

	roles := make([]*models.Role, 0)

	for _, key := range ul {
		u, err := s.GetRole(ctx, s.keyToCode(key))
		if err != nil {
			continue
		}

		roles = append(roles, &models.Role{
			ID:          u.ID,
			Code:        u.Code,
			Name:        u.Name,
			Description: u.Description,
			Blocked:     u.Blocked,
		})
	}

	return roles, nil
}

func (s *Roles) CreateRole(ctx context.Context, role models.RoleCreated) (*models.Role, error) {
	const op = "DbRoles.CreateRole"

	roleID := int(time.Now().Unix())
	roleCode := fmt.Sprintf("r%d", roleID)

	if _, err := s.client.HMSet(ctx, s.codeToKey(roleCode), redisclient.RoleCreated{
		ID:          roleID,
		Code:        roleCode,
		Name:        role.Name,
		Description: role.Description,
		Blocked:     role.Blocked,
	}).Result(); err != nil {
		return nil, fmt.Errorf("failed to create role | %s:%w", op, err)
	}

	return s.GetRole(ctx, roleCode)
}

func (s *Roles) UpdateRole(ctx context.Context, role models.RoleUpdated) (*models.Role, error) {
	const op = "DbRoles.UpdateRole"

	if _, err := s.GetRole(ctx, role.Code); err != nil {
		return nil, fmt.Errorf("failed to update role | %s:%w", op, err)
	}

	if _, err := s.client.HMSet(ctx, s.codeToKey(role.Code), redisclient.RoleUpdated{
		Name:        role.Name,
		Description: role.Description,
		Blocked:     role.Blocked,
	}).Result(); err != nil {
		return nil, fmt.Errorf("failed to update role | %s:%w", op, err)
	}

	return s.GetRole(ctx, role.Code)
}

func (s *Roles) DeleteRole(ctx context.Context, code string) error {
	const op = "DbRoles.DeleteRole"

	if _, err := s.GetRole(ctx, code); err != nil {
		return fmt.Errorf("failed to delete role | %s:%w", op, err)
	}

	if _, err := s.client.Del(ctx, s.codeToKey(code)).Result(); err != nil {
		return fmt.Errorf("failed to delete role | %s:%w", op, err)
	}

	return nil
}

func (s *Roles) space() string {
	return fmt.Sprintf("%s:*", space)
}

func (s *Roles) codeToKey(code string) string {
	return fmt.Sprintf("%s:%s", space, code)
}

func (s *Roles) keyToCode(key string) string {
	p := fmt.Sprintf("%s:", space)

	return strings.ReplaceAll(key, p, "")
}
