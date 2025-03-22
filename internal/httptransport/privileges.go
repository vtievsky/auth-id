package httptransport

import (
	"context"

	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) GetPrivileges(
	ctx context.Context,
	request serverhttp.GetPrivilegesRequestObject,
) (serverhttp.GetPrivilegesResponseObject, error) {
	roles, err := t.services.PrivilegeSvc.GetPrivileges(ctx)
	if err != nil {
		return serverhttp.GetPrivileges500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	resp := make([]serverhttp.Privilege, 0, len(roles))

	for _, role := range roles {
		resp = append(resp, serverhttp.Privilege{
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
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
