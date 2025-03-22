package httptransport

import (
	"context"

	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) GetPrivileges(
	ctx context.Context,
	request serverhttp.GetPrivilegesRequestObject,
) (serverhttp.GetPrivilegesResponseObject, error) {
	privileges, err := t.services.PrivilegeSvc.GetPrivileges(ctx, request.Params.PageSize, request.Params.Offset)
	if err != nil {
		return serverhttp.GetPrivileges500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	resp := make([]serverhttp.Privilege, 0, len(privileges))

	for _, privilege := range privileges {
		resp = append(resp, serverhttp.Privilege{
			Code:        privilege.Code,
			Name:        privilege.Name,
			Description: privilege.Description,
		})
	}

	return serverhttp.GetPrivileges200JSONResponse{
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}
