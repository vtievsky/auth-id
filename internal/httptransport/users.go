package httptransport

import (
	"context"

	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
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
	user, err := t.services.UserSvc.CreateUser(ctx, usersvc.UserCreated{
		Login:    request.Body.Login,
		FullName: request.Body.Name,
		Blocked:  request.Body.Blocked,
	})
	if err != nil {
		return serverhttp.CreateUser500JSONResponse{
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.CreateUser200JSONResponse{
		Data: serverhttp.User{
			Id:      user.ID,
			Login:   user.Login,
			Name:    user.FullName,
			Blocked: user.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) UpdateUser(ctx context.Context, request serverhttp.UpdateUserRequestObject) (serverhttp.UpdateUserResponseObject, error) {
	user, err := t.services.UserSvc.UpdateUser(ctx, usersvc.UserUpdated{
		Login:    request.Login,
		FullName: request.Body.Name,
		Blocked:  request.Body.Blocked,
	})
	if err != nil {
		return serverhttp.UpdateUser500JSONResponse{
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.UpdateUser200JSONResponse{
		Data: serverhttp.User{
			Id:      user.ID,
			Login:   user.Login,
			Name:    user.FullName,
			Blocked: user.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) DeleteUser(ctx context.Context, request serverhttp.DeleteUserRequestObject) (serverhttp.DeleteUserResponseObject, error) {
	if err := t.services.UserSvc.DeleteUser(ctx, request.Login); err != nil {
		return serverhttp.DeleteUser500JSONResponse{
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.DeleteUser200JSONResponse{
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}
