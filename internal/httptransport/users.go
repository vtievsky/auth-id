package httptransport

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
)

func (t *Transport) GetUser(
	ctx echo.Context,
	login string,
) error {
	user, err := t.services.UserSvc.GetUser(ctx.Request().Context(), login)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.GetUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.GetUserResponse200{ //nolint:wrapcheck
		Data: serverhttp.User{
			Name:    user.Name,
			Login:   user.Login,
			Blocked: user.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) GetUsers(ctx echo.Context, params serverhttp.GetUsersParams) error {
	users, err := t.services.UserSvc.GetUsers(
		ctx.Request().Context(),
		params.PageSize,
		params.Offset,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.GetUsersResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	resp := make([]serverhttp.User, 0, len(users))

	for _, user := range users {
		resp = append(resp, serverhttp.User{
			Name:    user.Name,
			Login:   user.Login,
			Blocked: user.Blocked,
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.GetUsersResponse200{ //nolint:wrapcheck
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) CreateUser(ctx echo.Context) error {
	var request serverhttp.CreateUserJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.CreateUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	user, err := t.services.UserSvc.CreateUser(ctx.Request().Context(), usersvc.UserCreated{
		Name:     request.Name,
		Login:    request.Login,
		Password: request.Password,
		Blocked:  request.Blocked,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.CreateUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.CreateUserResponse200{ //nolint:wrapcheck
		Data: serverhttp.User{
			Name:    user.Name,
			Login:   user.Login,
			Blocked: user.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) UpdateUser(
	ctx echo.Context,
	login string,
) error {
	var request serverhttp.UpdateUserJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.UpdateUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	// Запрет блокирования самого себя
	if request.Blocked {
		if err := t.yourSelf(ctx, login); err != nil {
			err = fmt.Errorf("failed to update user | %w", ErrBlockHimself)

			return ctx.JSON(http.StatusInternalServerError, serverhttp.DeleteUserResponse500{ //nolint:wrapcheck
				Status: serverhttp.ResponseStatusError{
					Code:        serverhttp.Error,
					Description: err.Error(),
				},
			})
		}
	}

	user, err := t.services.UserSvc.UpdateUser(ctx.Request().Context(), usersvc.UserUpdated{
		Name:    request.Name,
		Login:   login,
		Blocked: request.Blocked,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.UpdateUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.UpdateUserResponse200{ //nolint:wrapcheck
		Data: serverhttp.User{
			Name:    user.Name,
			Login:   user.Login,
			Blocked: user.Blocked,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) DeleteUser(
	ctx echo.Context,
	login string,
) error {
	// Запрет удаления самого себя
	if err := t.yourSelf(ctx, login); err != nil {
		err = fmt.Errorf("failed to delete user | %w", ErrDeleteHimself)

		return ctx.JSON(http.StatusInternalServerError, serverhttp.DeleteUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	if err := t.services.UserSvc.DeleteUser(ctx.Request().Context(), login); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.DeleteUserResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.DeleteUserResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) ChangePass(
	ctx echo.Context,
	login string,
) error {
	var request serverhttp.ChangePassJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.ChangePassResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	if err := t.services.UserSvc.ChangePass(ctx.Request().Context(), login, request.Current, request.Changed); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.ChangePassResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.ChangePassResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) ResetPass(
	ctx echo.Context,
	login string,
) error {
	var request serverhttp.ResetPassJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.ResetPassResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	if err := t.services.UserSvc.ResetPass(ctx.Request().Context(), login, request.Changed); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.ResetPassResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.ResetPassResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}
