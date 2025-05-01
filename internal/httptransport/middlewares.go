package httptransport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	oteltracing "github.com/vtievsky/auth-id/internal/otel/tracing"
	authidjwt "github.com/vtievsky/auth-id/pkg/jwt"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type SessionService interface {
	Search(ctx context.Context, sessionID, privilegeCode string) error
}

func LoggerMiddleware(l *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, span := otel.Tracer("").Start(c.Request().Context(), "login")
			defer span.End()

			c.SetRequest(c.Request().WithContext(ctx))

			startTime := time.Now()

			err := next(c)

			statusCode := c.Response().Status

			if errors.Is(c.Request().Context().Err(), context.Canceled) {
				statusCode = 499 // Client CLosed Request
			}

			if err != nil || statusCode >= http.StatusInternalServerError {
				l.Error("request",
					zap.Error(err),
					zap.String("ip", c.RealIP()),
					zap.String("trace_id", oteltracing.GetTraceID(ctx)),
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
				zap.String("trace_id", oteltracing.GetTraceID(ctx)),
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
	sessionSvc SessionService,
	signingKey string,
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
			endpointPrivilegeKey := endpointPrivilegeKey(c)

			if _, ok := endpointWithout[endpointPrivilegeKey]; ok {
				return next(c)
			}

			endpointPrivilegeCode, ok := endpointWithPrivileges[endpointPrivilegeKey]

			if !ok {
				return fmt.Errorf("privilege path could not be mapped")
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

			err = sessionSvc.Search(c.Request().Context(), token.SessionID, endpointPrivilegeCode)
			if err != nil {
				return fmt.Errorf("failed to search session privilege | %w", err)
			}

			c.Set("session_id", token.SessionID)

			return next(c)
		}
	}
}

func endpointPrivilegeKey(c echo.Context) string {
	return fmt.Sprintf("%s%s", strings.ToLower(c.Request().Method), c.Path())
}
