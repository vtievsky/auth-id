package httptransport

import (
	"context"

	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
)

func (t *Transport) GetRolePrivileges(
	ctx context.Context,
	request serverhttp.GetRolePrivilegesRequestObject,
) (serverhttp.GetRolePrivilegesResponseObject, error) {
	privileges, err := t.services.RoleSvc.GetRolePrivileges(ctx, request.Code)
	if err != nil {
		return serverhttp.GetRolePrivileges500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	resp := make([]serverhttp.RolePrivilege, 0, len(privileges))

	for _, role := range privileges {
		resp = append(resp, serverhttp.RolePrivilege{
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			Allowed:     role.Allowed,
		})
	}

	return serverhttp.GetRolePrivileges200JSONResponse{
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) AddRolePrivilege(
	ctx context.Context,
	request serverhttp.AddRolePrivilegeRequestObject,
) (serverhttp.AddRolePrivilegeResponseObject, error) {
	if err := t.services.RoleSvc.AddRolePrivilege(ctx, rolesvc.RolePrivilegeCreated{
		RoleCode:      request.RoleCode,
		PrivilegeCode: request.PrivilegeCode,
		Allowed:       request.Body.Allowed,
	}); err != nil {
		return serverhttp.AddRolePrivilege500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.AddRolePrivilege200JSONResponse{
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) UpdateRolePrivilege(
	ctx context.Context,
	request serverhttp.UpdateRolePrivilegeRequestObject,
) (serverhttp.UpdateRolePrivilegeResponseObject, error) {
	if err := t.services.RoleSvc.UpdateRolePrivilege(ctx, rolesvc.RolePrivilegeUpdated{
		RoleCode:      request.RoleCode,
		PrivilegeCode: request.PrivilegeCode,
		Allowed:       request.Body.Allowed,
	}); err != nil {
		return serverhttp.UpdateRolePrivilege500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.UpdateRolePrivilege200JSONResponse{
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) DeleteRolePrivilege(
	ctx context.Context,
	request serverhttp.DeleteRolePrivilegeRequestObject,
) (serverhttp.DeleteRolePrivilegeResponseObject, error) {
	if err := t.services.RoleSvc.DeleteRolePrivilege(ctx, rolesvc.RolePrivilegeDeleted{
		RoleCode:      request.RoleCode,
		PrivilegeCode: request.PrivilegeCode,
	}); err != nil {
		return serverhttp.DeleteRolePrivilege500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.DeleteRolePrivilege200JSONResponse{
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}
