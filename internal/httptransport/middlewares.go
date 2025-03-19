package httptransport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	authidjwt "github.com/vtievsky/auth-id/pkg/jwt"
	"go.uber.org/zap"
)

func LoggerMiddleware(l *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			startTime := time.Now()
			statusCode := c.Response().Status

			if errors.Is(c.Request().Context().Err(), context.Canceled) {
				statusCode = 499 // Client CLosed Request
			}

			if err != nil || statusCode >= http.StatusInternalServerError {
				l.Error("request",
					zap.Error(err),
					zap.String("ip", c.RealIP()),
					// zap.String("trace_id", motel.GetTraceID(c.Request().Context())),
					zap.String("method", c.Request().Method),
					zap.String("path", c.Request().RequestURI),
					zap.String("host", c.Request().Host),
					zap.String("duration", time.Since(startTime).String()),
					zap.Int("status_code", statusCode),
				)

				return err
			}

			l.Info("request",
				zap.String("ip", c.RealIP()),
				// zap.String("trace_id", motel.GetTraceID(c.Request().Context())),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().RequestURI),
				zap.String("host", c.Request().Host),
				zap.String("duration", time.Since(startTime).String()),
				zap.Int("status_code", statusCode),
			)

			return nil
		}
	}
}

func AuthorizationMiddleware(
	signingKey string,
	searchPrivilegeFunc func(ctx context.Context, sessionID, privilegeCode string) error,
) echo.MiddlewareFunc {
	extractTokenValue := func(header http.Header) (string, error) {
		values, ok := header["Authorization"]
		if !ok {
			return "", fmt.Errorf("token not found")
		}

		if len(values) < 1 {
			return "", fmt.Errorf("token not found")
		}

		ul := strings.Split(values[0], " ")

		if len(ul) < 2 { //nolint:mnd
			return "", fmt.Errorf("invalid token")
		}

		return ul[1], nil
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().RequestURI == "/v1/users/slaventius/sessions" {
				return next(c)
			}

			value, err := extractTokenValue(c.Request().Header)
			if err != nil {
				return err
			}

			token, err := authidjwt.ParseToken([]byte(signingKey), []byte(value))
			if err != nil {
				return err //nolint:wrapcheck
			}

			if !token.Valid {
				return fmt.Errorf("token not valid")
			}

			err = searchPrivilegeFunc(c.Request().Context(), token.SessionID, "user_read")
			if err != nil {
				return err
			}

			return next(c)
		}
	}
}
