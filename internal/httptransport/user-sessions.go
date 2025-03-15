package httptransport

import (
	"context"

	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
)

func (t *Transport) Login(
	ctx context.Context,
	request serverhttp.LoginRequestObject,
) (serverhttp.LoginResponseObject, error) {
	session, err := t.services.SessionSvc.Login(ctx, request.Login, request.Body.Password)
	if err != nil {
		return serverhttp.Login500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.Login200JSONResponse{
		Data: serverhttp.Session{
			Id:        string(session.ID),
			CreatedAt: session.CreatedAt,
			ExpiredAt: session.ExpiredAt,
		},
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (t *Transport) GetUserSessions(
	ctx context.Context,
	request serverhttp.GetUserSessionsRequestObject,
) (serverhttp.GetUserSessionsResponseObject, error) {
	sessions, err := t.services.SessionSvc.GetUserSessions(ctx, request.Login)
	if err != nil {
		return serverhttp.GetUserSessions500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	resp := make([]serverhttp.Session, 0, len(sessions))

	for _, session := range sessions {
		resp = append(resp, serverhttp.Session{
			Id:        string(session.ID),
			CreatedAt: session.CreatedAt,
			ExpiredAt: session.ExpiredAt,
		})
	}

	return serverhttp.GetUserSessions200JSONResponse{
		Data: resp,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}
