package httptransport

import (
	"context"
	"time"

	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
)

func (t *Transport) GetRoleUsers(
	ctx context.Context,
	request serverhttp.GetRoleUsersRequestObject,
) (serverhttp.GetRoleUsersResponseObject, error) {
	users, err := t.services.RoleUserSvc.GetRoleUsers(ctx, request.Code)
	if err != nil {
		return serverhttp.GetRoleUsers500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	resp := make([]serverhttp.RoleUser, 0, len(users))

	for _, user := range users {
		resp = append(resp, serverhttp.RoleUser{
			Name:    user.Name,
			Login:   user.Login,
			DateIn:  user.DateIn.Format(time.RFC3339),
			DateOut: user.DateOut.Format(time.RFC3339),
		})
	}

	return serverhttp.GetRoleUsers200JSONResponse{
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) AddRoleUser(
	ctx context.Context,
	request serverhttp.AddRoleUserRequestObject,
) (serverhttp.AddRoleUserResponseObject, error) {
	if err := t.services.RoleUserSvc.AddRoleUser(ctx, rolesvc.RoleUserCreated{
		Login:    request.Login,
		RoleCode: request.RoleCode,
		DateIn:   time.Now(),                               // TODO request.Body.DateIn,
		DateOut:  time.Now().Add(time.Hour * 24 * 30 * 12), // TODO request.Body.DateOut,
	}); err != nil {
		return serverhttp.AddRoleUser500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.AddRoleUser200JSONResponse{
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

// func (t *Transport) UpdateRoleUser(
// 	ctx context.Context,
// 	request serverhttp.UpdateRoleUserRequestObject,
// ) (serverhttp.UpdateRoleUserResponseObject, error) {
// 	if err := t.services.RoleUserSvc.UpdateRoleUser(ctx, rolesvc.RoleUserUpdated{
// 		RoleCode:      request.RoleCode,
// 		PrivilegeCode: request.PrivilegeCode,
// 		Allowed:       request.Body.Allowed,
// 	}); err != nil {
// 		return serverhttp.UpdateRoleUser500JSONResponse{ //nolint:nilerr
// 			Status: serverhttp.ResponseStatusError{
// 				Code:        serverhttp.Error,
// 				Description: err.Error(),
// 			},
// 		}, nil
// 	}

// 	return serverhttp.UpdateRoleUser200JSONResponse{
// 		Status: serverhttp.ResponseStatusOk{
// 			Code:        serverhttp.Ok,
// 			Description: "",
// 		},
// 	}, nil
// }

// func (t *Transport) DeleteRoleUser(
// 	ctx context.Context,
// 	request serverhttp.DeleteRoleUserRequestObject,
// ) (serverhttp.DeleteRoleUserResponseObject, error) {
// 	if err := t.services.RoleUserSvc.DeleteRoleUser(ctx, rolesvc.RoleUserDeleted{
// 		RoleCode:      request.RoleCode,
// 		PrivilegeCode: request.PrivilegeCode,
// 	}); err != nil {
// 		return serverhttp.DeleteRoleUser500JSONResponse{ //nolint:nilerr
// 			Status: serverhttp.ResponseStatusError{
// 				Code:        serverhttp.Error,
// 				Description: err.Error(),
// 			},
// 		}, nil
// 	}

// 	return serverhttp.DeleteRoleUser200JSONResponse{
// 		Status: serverhttp.ResponseStatusOk{
// 			Code:        serverhttp.Ok,
// 			Description: "",
// 		},
// 	}, nil
// }
