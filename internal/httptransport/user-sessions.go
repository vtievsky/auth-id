package httptransport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) Login(
	ctx echo.Context,
	login string,
) error {
	var request serverhttp.LoginJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.LoginResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	resp, err := t.services.SessionSvc.Login(ctx.Request().Context(), login, request.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.LoginResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.LoginResponse200{ //nolint:wrapcheck
		Data: serverhttp.ResponseAccess{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) GetUserSessions(
	ctx echo.Context,
	login string,
	params serverhttp.GetUserSessionsParams,
) error {
	sessions, err := t.services.SessionSvc.GetUserSessions(
		ctx.Request().Context(),
		login,
		params.PageSize,
		params.Offset,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.GetUserSessionsResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	resp := make([]serverhttp.Session, 0, len(sessions))

	for _, session := range sessions {
		resp = append(resp, serverhttp.Session{
			Id:        string(session.ID),
			CreatedAt: session.CreatedAt,
			ExpiredAt: session.ExpiredAt,
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.GetUserSessionsResponse200{ //nolint:wrapcheck
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}

func (t *Transport) DeleteUserSession(
	ctx echo.Context,
	login, sessionID string,
) error {
	if err := t.services.SessionSvc.Delete(ctx.Request().Context(), login, sessionID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, serverhttp.DeleteUserSessionResponse500{ //nolint:wrapcheck
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, serverhttp.DeleteUserSessionResponse200{ //nolint:wrapcheck
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	})
}
