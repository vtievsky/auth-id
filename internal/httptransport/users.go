package httptransport

import (
	"context"

	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) GetUsers(ctx context.Context, request serverhttp.GetUsersRequestObject) (serverhttp.GetUsersResponseObject, error) {
	return serverhttp.GetUsers200JSONResponse{
		Data: []serverhttp.User{},
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
