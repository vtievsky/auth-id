package httptransport

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/vtievsky/auth-id/internal/conf"
	"github.com/vtievsky/auth-id/internal/services"
)

type Transport struct {
	conf     *conf.Config
	services *services.SvcLayer
}

func New(conf *conf.Config, svc *services.SvcLayer) *Transport {
	return &Transport{
		conf:     conf,
		services: svc,
	}
}

func (t *Transport) yourSelf(ctx echo.Context, login string) error {
	sessionID, ok := ctx.Get("session_id").(string)
	if !ok {
		return fmt.Errorf("failed to assert type session id")
	}

	cart, err := t.services.SessionSvc.Get(ctx.Request().Context(), sessionID)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if strings.EqualFold(login, cart.Login) {
		return ErrHimself
	}

	return nil
}
