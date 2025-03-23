package httptransport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) GetUserPrivileges(
	ctx echo.Context,
	login string,
	params serverhttp.GetUserPrivilegesParams,
) error {
	privileges, err := t.services.UserPrivilegeSvc.GetUserPrivileges(
		ctx.Request().Context(),
		login,
		params.PageSize,
		params.Offset,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.GetUserPrivilegesResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	resp := make([]serverhttp.UserPrivilege, 0, len(privileges))

	for _, privilege := range privileges {
		resp = append(resp, serverhttp.UserPrivilege{
			Code:        privilege.Code,
			Name:        privilege.Name,
			Description: privilege.Description,
			DateIn:      types.Date{Time: privilege.DateIn},
			DateOut:     types.Date{Time: privilege.DateOut},
		})
	}

	return ctx.JSON(http.StatusInternalServerError, serverhttp.GetUserPrivilegesResponse200{ //nolint:wrapcheck
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}
