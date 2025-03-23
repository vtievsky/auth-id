package httptransport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) GetPrivileges(ctx echo.Context, params serverhttp.GetPrivilegesParams) error {
	privileges, err := t.services.PrivilegeSvc.GetPrivileges(
		ctx.Request().Context(),
		params.PageSize,
		params.Offset,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.GetPrivilegesResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	resp := make([]serverhttp.Privilege, 0, len(privileges))

	for _, privilege := range privileges {
		resp = append(resp, serverhttp.Privilege{
			Code:        privilege.Code,
			Name:        privilege.Name,
			Description: privilege.Description,
		})
	}

	return ctx.JSON(http.StatusInternalServerError, serverhttp.GetPrivilegesResponse200{ //nolint:wrapcheck
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}
