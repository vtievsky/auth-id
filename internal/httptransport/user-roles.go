package httptransport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) GetUserRoles(
	ctx echo.Context,
	login string,
	params serverhttp.GetUserRolesParams,
) error {
	roles, err := t.services.UserRoleSvc.GetUserRoles(
		ctx.Request().Context(),
		login,
		params.PageSize,
		params.Offset,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.GetUserRolesResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	resp := make([]serverhttp.UserRole, 0, len(roles))

	for _, role := range roles {
		resp = append(resp, serverhttp.UserRole{
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			DateIn:      types.Date{Time: role.DateIn},
			DateOut:     types.Date{Time: role.DateOut},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.GetUserRolesResponse200{ //nolint:wrapcheck
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}
