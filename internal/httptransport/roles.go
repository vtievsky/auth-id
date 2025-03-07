package httptransport

import (
	"context"

	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
)

func (t *Transport) GetRole(
	ctx context.Context,
	request serverhttp.GetRoleRequestObject,
) (serverhttp.GetRoleResponseObject, error) {
	role, err := t.services.RoleSvc.GetRole(ctx, request.Id)
	if err != nil {
		return serverhttp.GetRole500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.GetRole200JSONResponse{
		Data: serverhttp.Role{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Blocked:     role.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) GetRoles(ctx context.Context, request serverhttp.GetRolesRequestObject) (serverhttp.GetRolesResponseObject, error) {
	roles, err := t.services.RoleSvc.GetRoles(ctx)
	if err != nil {
		return serverhttp.GetRoles500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	resp := make([]serverhttp.Role, 0, len(roles))

	for _, role := range roles {
		resp = append(resp, serverhttp.Role{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Blocked:     role.Blocked,
		})
	}

	return serverhttp.GetRoles200JSONResponse{
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) CreateRole(
	ctx context.Context,
	request serverhttp.CreateRoleRequestObject,
) (serverhttp.CreateRoleResponseObject, error) {
	role, err := t.services.RoleSvc.CreateRole(ctx, rolesvc.RoleCreated{
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Blocked:     request.Body.Blocked,
	})
	if err != nil {
		return serverhttp.CreateRole500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.CreateRole200JSONResponse{
		Data: serverhttp.Role{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Blocked:     role.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) UpdateRole(
	ctx context.Context,
	request serverhttp.UpdateRoleRequestObject,
) (serverhttp.UpdateRoleResponseObject, error) {
	role, err := t.services.RoleSvc.UpdateRole(ctx, rolesvc.RoleUpdated{
		ID:          request.Id,
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Blocked:     request.Body.Blocked,
	})
	if err != nil {
		return serverhttp.UpdateRole500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.UpdateRole200JSONResponse{
		Data: serverhttp.Role{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Blocked:     role.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) DeleteRole(
	ctx context.Context,
	request serverhttp.DeleteRoleRequestObject,
) (serverhttp.DeleteRoleResponseObject, error) {
	if err := t.services.RoleSvc.DeleteRole(ctx, request.Id); err != nil {
		return serverhttp.DeleteRole500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.DeleteRole200JSONResponse{
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}
