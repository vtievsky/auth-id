package httptransport

import (
	"context"

	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
)

func (t *Transport) GetUser(
	ctx context.Context,
	request serverhttp.GetUserRequestObject,
) (serverhttp.GetUserResponseObject, error) {
	user, err := t.services.UserSvc.GetUser(ctx, request.Login)
	if err != nil {
		return serverhttp.GetUser500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.GetUser200JSONResponse{
		Data: serverhttp.User{
			Name:    user.Name,
			Login:   user.Login,
			Blocked: user.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) GetUsers(ctx context.Context, request serverhttp.GetUsersRequestObject) (serverhttp.GetUsersResponseObject, error) {
	users, err := t.services.UserSvc.GetUsers(ctx)
	if err != nil {
		return serverhttp.GetUsers500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	resp := make([]serverhttp.User, 0, len(users))

	for _, user := range users {
		resp = append(resp, serverhttp.User{
			Name:    user.Name,
			Login:   user.Login,
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

func (t *Transport) CreateUser(
	ctx context.Context,
	request serverhttp.CreateUserRequestObject,
) (serverhttp.CreateUserResponseObject, error) {
	user, err := t.services.UserSvc.CreateUser(ctx, usersvc.UserCreated{
		Name:     request.Body.Name,
		Login:    request.Body.Login,
		Password: request.Body.Password,
		Blocked:  request.Body.Blocked,
	})
	if err != nil {
		return serverhttp.CreateUser500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.CreateUser200JSONResponse{
		Data: serverhttp.User{
			Name:    user.Name,
			Login:   user.Login,
			Blocked: user.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) UpdateUser(
	ctx context.Context,
	request serverhttp.UpdateUserRequestObject,
) (serverhttp.UpdateUserResponseObject, error) {
	user, err := t.services.UserSvc.UpdateUser(ctx, usersvc.UserUpdated{
		Name:    request.Body.Name,
		Login:   request.Login,
		Blocked: request.Body.Blocked,
	})
	if err != nil {
		return serverhttp.UpdateUser500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.UpdateUser200JSONResponse{
		Data: serverhttp.User{
			Name:    user.Name,
			Login:   user.Login,
			Blocked: user.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) DeleteUser(
	ctx context.Context,
	request serverhttp.DeleteUserRequestObject,
) (serverhttp.DeleteUserResponseObject, error) {
	if err := t.services.UserSvc.DeleteUser(ctx, request.Login); err != nil {
		return serverhttp.DeleteUser500JSONResponse{ //nolint:nilerr
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
