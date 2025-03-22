package userprivilegesvc

import (
	"context"
	"fmt"
	"math"
	"time"

	roleprivilegesvc "github.com/vtievsky/auth-id/internal/services/role-privileges"
	userrolesvc "github.com/vtievsky/auth-id/internal/services/user-roles"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	threadsLimit = 10
)

type userPrivilegeStats struct {
	Code        string
	Name        string
	Description string
	Allowed     bool
	DateIn      time.Time
	DateOut     time.Time
}

type UserPrivilege struct {
	Code        string
	Name        string
	Description string
	DateIn      time.Time
	DateOut     time.Time
}

type UserRoleSvc interface {
	GetUserRoles(ctx context.Context, login string, pageSize, offset uint32) ([]*userrolesvc.UserRole, error)
}

type RolePrivilegeSvc interface {
	GetRolePrivileges(ctx context.Context, code string, pageSize, offset uint32) ([]*roleprivilegesvc.RolePrivilege, error)
}

type UserPrivilegeSvcOpts struct {
	Logger           *zap.Logger
	UserRoleSvc      UserRoleSvc
	RolePrivilegeSvc RolePrivilegeSvc
}

type UserPrivilegeSvc struct {
	logger           *zap.Logger
	userRoleSvc      UserRoleSvc
	rolePrivilegeSvc RolePrivilegeSvc
}

func New(opts *UserPrivilegeSvcOpts) *UserPrivilegeSvc {
	return &UserPrivilegeSvc{
		logger:           opts.Logger,
		userRoleSvc:      opts.UserRoleSvc,
		rolePrivilegeSvc: opts.RolePrivilegeSvc,
	}
}

func (s *UserPrivilegeSvc) GetUserPrivileges(ctx context.Context, login string, pageSize, offset uint32) ([]*UserPrivilege, error) {
	const op = "UserPrivilegeSvc.GetUserPrivileges"

	fetchRolePrivileges := func(actx context.Context, acombine chan<- userPrivilegeStats,
		aroleCode string, adateIn, adateOut time.Time) func() error {
		return func() error {
			rolePrivileges, err := s.rolePrivilegeSvc.GetRolePrivileges(actx, aroleCode, math.MaxUint32, 0)
			if err != nil {
				return fmt.Errorf("failed to fetch role privileges %s:%w", op, err)
			}

			for _, privilege := range rolePrivileges {
				acombine <- userPrivilegeStats{
					Code:        privilege.Code,
					Name:        privilege.Name,
					Description: privilege.Description,
					Allowed:     privilege.Allowed,
					DateIn:      adateIn,
					DateOut:     adateOut,
				}
			}

			return nil
		}
	}

	userRoles, err := s.userRoleSvc.GetUserRoles(ctx, login, pageSize, offset)
	if err != nil {
		s.logger.Error("failed to get user roles",
			zap.String("login", login),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to get user roles | %s:%w", op, err)
	}

	combineStats := make(chan userPrivilegeStats)

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(threadsLimit + 1)

	g.Go(func() error {
		for _, role := range userRoles {
			g.Go(fetchRolePrivileges(gCtx, combineStats, role.Code, role.DateIn, role.DateOut))
		}

		return nil
	})

	go func() {
		err = g.Wait()

		close(combineStats)
	}()

	var (
		ok        bool
		privilege *userPrivilegeStats
	)

	mapRolePrivilege := make(map[string]*userPrivilegeStats)

	for v := range combineStats {
		if privilege, ok = mapRolePrivilege[v.Code]; ok {
			// Увеличиваем период действия привилегии
			if v.DateIn.Before(privilege.DateIn) {
				privilege.DateIn = v.DateIn
			}

			if privilege.DateOut.Before(v.DateOut) {
				privilege.DateOut = v.DateOut
			}

			// Запрет имеет большее преимущество
			privilege.Allowed = privilege.Allowed && v.Allowed

			continue
		}

		mapRolePrivilege[v.Code] = &userPrivilegeStats{
			Code:        v.Code,
			Name:        v.Name,
			Description: v.Description,
			Allowed:     v.Allowed,
			DateIn:      v.DateIn,
			DateOut:     v.DateOut,
		}
	}

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	resp := make([]*UserPrivilege, 0)

	for _, privilege := range mapRolePrivilege {
		if !privilege.Allowed {
			continue
		}

		resp = append(resp, &UserPrivilege{
			Code:        privilege.Code,
			Name:        privilege.Name,
			Description: privilege.Description,
			DateIn:      privilege.DateIn,
			DateOut:     privilege.DateOut,
		})
	}

	return resp, nil
}
