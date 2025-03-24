package httptransport

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
	roleusersvc "github.com/vtievsky/auth-id/internal/services/role-users"
)

func (t *Transport) GetRoleUsers(
	ctx echo.Context,
	code string,
	params serverhttp.GetRoleUsersParams,
) error {
	users, err := t.services.RoleUserSvc.GetRoleUsers(
		ctx.Request().Context(),
		code,
		params.PageSize,
		params.Offset,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.GetRoleUsersResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	resp := make([]serverhttp.RoleUser, 0, len(users))

	for _, user := range users {
		resp = append(resp, serverhttp.RoleUser{
			Name:    user.Name,
			Login:   user.Login,
			DateIn:  types.Date{Time: user.DateIn},
			DateOut: types.Date{Time: user.DateOut},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.GetRoleUsersResponse200{ //nolint:wrapcheck
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) AddRoleUser(
	ctx echo.Context,
	roleCode, login string,
) error {
	// Запрет добавления собственной роли
	if err := t.yourSelf(ctx, login); err != nil {
		err = fmt.Errorf("failed to add role user | %w", ErrAddHimselfRoles)

		return ctx.JSON(http.StatusInternalServerError, serverhttp.AddRoleUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	var request serverhttp.AddRoleUserJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.AddRoleUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	if err := t.services.RoleUserSvc.AddRoleUser(ctx.Request().Context(),
		roleusersvc.RoleUserCreated{
			Login:    login,
			RoleCode: roleCode,
			DateIn:   request.DateIn.Time,
			DateOut:  request.DateOut.Time,
		},
	); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.AddRoleUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.AddRoleUserResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) UpdateRoleUser(
	ctx echo.Context,
	roleCode, login string,
) error {
	var request serverhttp.UpdateRoleUserJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.AddRoleUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	if err := t.services.RoleUserSvc.UpdateRoleUser(ctx.Request().Context(),
		roleusersvc.RoleUserUpdated{
			Login:    login,
			RoleCode: roleCode,
			DateIn:   request.DateIn.Time,
			DateOut:  request.DateOut.Time,
		},
	); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.UpdateRoleUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.UpdateRoleUserResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) DeleteRoleUser(
	ctx echo.Context,
	roleCode, login string,
) error {
	// Запрет удаления собственной роли
	if err := t.yourSelf(ctx, login); err != nil {
		err = fmt.Errorf("failed to delete role user | %w", ErrDeleteHimselfRoles)

		return ctx.JSON(http.StatusInternalServerError, serverhttp.DeleteRoleUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	if err := t.services.RoleUserSvc.DeleteRoleUser(ctx.Request().Context(),
		roleusersvc.RoleUserDeleted{
			Login:    login,
			RoleCode: roleCode,
		},
	); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.DeleteRoleUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.DeleteRoleUserResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}
