package httptransport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
	roleprivilegesvc "github.com/vtievsky/auth-id/internal/services/role-privileges"
)

func (t *Transport) GetRolePrivileges(
	ctx echo.Context,
	code string,
	params serverhttp.GetRolePrivilegesParams,
) error {
	privileges, err := t.services.RolePrivilegeSvc.GetRolePrivileges(
		ctx.Request().Context(),
		code,
		params.PageSize,
		params.Offset,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.GetRolePrivilegesResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
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

	return ctx.JSON(http.StatusOK, serverhttp.GetRolePrivilegesResponse200{ //nolint:wrapcheck
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) AddRolePrivilege(
	ctx echo.Context,
	roleCode, privilegeCode string,
) error {
	var request serverhttp.AddRolePrivilegeJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.AddRolePrivilegeResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	if err := t.services.RolePrivilegeSvc.AddRolePrivilege(ctx.Request().Context(),
		roleprivilegesvc.RolePrivilegeCreated{
			RoleCode:      roleCode,
			PrivilegeCode: privilegeCode,
			Allowed:       request.Allowed,
		},
	); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.AddRolePrivilegeResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.AddRolePrivilegeResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) UpdateRolePrivilege(
	ctx echo.Context,
	roleCode, privilegeCode string,
) error {
	var request serverhttp.UpdateRolePrivilegeJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.AddRolePrivilegeResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	if err := t.services.RolePrivilegeSvc.UpdateRolePrivilege(ctx.Request().Context(),
		roleprivilegesvc.RolePrivilegeUpdated{
			RoleCode:      roleCode,
			PrivilegeCode: privilegeCode,
			Allowed:       request.Allowed,
		},
	); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.UpdateRolePrivilegeResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.UpdateRolePrivilegeResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) DeleteRolePrivilege(
	ctx echo.Context,
	roleCode, privilegeCode string,
) error {
	if err := t.services.RolePrivilegeSvc.DeleteRolePrivilege(ctx.Request().Context(),
		roleprivilegesvc.RolePrivilegeDeleted{
			RoleCode:      roleCode,
			PrivilegeCode: privilegeCode,
		},
	); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.DeleteRolePrivilegeResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.DeleteRolePrivilegeResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}
