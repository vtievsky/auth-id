package httptransport

import (
	"context"

	"github.com/oapi-codegen/runtime/types"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) GetUserRoles(
	ctx context.Context,
	request serverhttp.GetUserRolesRequestObject,
) (serverhttp.GetUserRolesResponseObject, error) {
	roles, err := t.services.UserRoleSvc.GetUserRoles(ctx, request.Login)
	if err != nil {
		return serverhttp.GetUserRoles500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
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

	return serverhttp.GetUserRoles200JSONResponse{
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}
