package httptransport

import (
	"context"

	"github.com/oapi-codegen/runtime/types"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) GetUserPrivileges(
	ctx context.Context,
	request serverhttp.GetUserPrivilegesRequestObject,
) (serverhttp.GetUserPrivilegesResponseObject, error) {
	privileges, err := t.services.UserPrivilegeSvc.GetUserPrivileges(ctx, request.Login)
	if err != nil {
		return serverhttp.GetUserPrivileges500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
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

	return serverhttp.GetUserPrivileges200JSONResponse{
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}
