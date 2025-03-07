package rolesvc

import (
	"context"

	"go.uber.org/zap"
)

type Role struct {
	ID          int
	Name        string
	Description string
	Blocked     bool
}

type RoleCreated struct {
	Name        string
	Description string
	Blocked     bool
}

type RoleUpdated struct {
	ID          int
	Name        string
	Description string
	Blocked     bool
}

type Storage interface {
}

type RoleSvcOpts struct {
	Logger  *zap.Logger
	Storage Storage
}

type RoleSvc struct {
	logger  *zap.Logger
	storage Storage
}

func New(opts *RoleSvcOpts) *RoleSvc {
	return &RoleSvc{
		logger:  opts.Logger,
		storage: opts.Storage,
	}
}

func (s *RoleSvc) GetRole(ctx context.Context, id int) (*Role, error) {
	const op = "RoleSvc.GetUser"

	return &Role{
		ID:          id,
		Name:        "",
		Description: "",
		Blocked:     false,
	}, nil
}

func (s *RoleSvc) GetRoles(ctx context.Context) ([]*Role, error) {
	const op = "RoleSvc.GetRoles"

	return []*Role{}, nil
}

func (s *RoleSvc) CreateRole(ctx context.Context, user RoleCreated) (*Role, error) {
	const op = "RoleSvc.CreateRole"

	return &Role{
		ID:          0,
		Name:        "",
		Description: "",
		Blocked:     false,
	}, nil
}

func (s *RoleSvc) UpdateRole(ctx context.Context, user RoleUpdated) (*Role, error) {
	const op = "RoleSvc.UpdateRole"

	return &Role{
		ID:          0,
		Name:        "",
		Description: "",
		Blocked:     false,
	}, nil
}

func (s *RoleSvc) DeleteRole(ctx context.Context, id int) error {
	const op = "RoleSvc.RoleUser"

	return nil
}
