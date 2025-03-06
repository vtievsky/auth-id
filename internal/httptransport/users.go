package httptransport

import (
	"context"

	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) GetUsers(ctx context.Context, request serverhttp.GetUsersRequestObject) (serverhttp.GetUsersResponseObject, error) {
	users, err := t.services.UserSvc.GetUsers(ctx)
	if err != nil {
		return serverhttp.GetUsers500JSONResponse{
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	resp := make([]serverhttp.User, 0, len(users))

	for _, user := range users {
		resp = append(resp, serverhttp.User{
			Id:      user.ID,
			Login:   user.Login,
			Name:    user.FullName,
			Blocked: user.Blocked,
		})
	}

	return serverhttp.GetUsers200JSONResponse{
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) CreateUser(ctx context.Context, request serverhttp.CreateUserRequestObject) (serverhttp.CreateUserResponseObject, error) {
	return serverhttp.CreateUser200JSONResponse{
		Data: serverhttp.User{},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}
