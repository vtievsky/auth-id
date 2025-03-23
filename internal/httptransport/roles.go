package httptransport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
)

func (t *Transport) GetRole(
	ctx echo.Context,
	code string,
) error {
	role, err := t.services.RoleSvc.GetRole(ctx.Request().Context(), code)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.GetRoleResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.GetRoleResponse200{ //nolint:wrapcheck
		Data: serverhttp.Role{
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			Blocked:     role.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) GetRoles(
	ctx echo.Context,
	params serverhttp.GetRolesParams,
) error {
	roles, err := t.services.RoleSvc.GetRoles(ctx.Request().Context(), params.PageSize, params.Offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.GetRolesResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	resp := make([]serverhttp.Role, 0, len(roles))

	for _, role := range roles {
		resp = append(resp, serverhttp.Role{
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			Blocked:     role.Blocked,
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.GetRolesResponse200{ //nolint:wrapcheck
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) CreateRole(ctx echo.Context) error {
	var request serverhttp.CreateRoleJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.CreateRoleResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	role, err := t.services.RoleSvc.CreateRole(ctx.Request().Context(),
		rolesvc.RoleCreated{
			Name:        request.Name,
			Description: request.Description,
			Blocked:     request.Blocked,
		},
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.CreateRoleResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.CreateRoleResponse200{ //nolint:wrapcheck
		Data: serverhttp.Role{
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			Blocked:     role.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) UpdateRole(
	ctx echo.Context,
	code string,
) error {
	var request serverhttp.UpdateRoleJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.UpdateRoleResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	role, err := t.services.RoleSvc.UpdateRole(ctx.Request().Context(),
		rolesvc.RoleUpdated{
			Code:        code,
			Name:        request.Name,
			Description: request.Description,
			Blocked:     request.Blocked,
		},
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.UpdateRoleResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.UpdateRoleResponse200{ //nolint:wrapcheck
		Data: serverhttp.Role{
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			Blocked:     role.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) DeleteRole(
	ctx echo.Context,
	code string,
) error {
	if err := t.services.RoleSvc.DeleteRole(ctx.Request().Context(), code); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.DeleteRoleResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.DeleteRoleResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}
