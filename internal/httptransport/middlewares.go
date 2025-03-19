package httptransport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
	"go.uber.org/zap"
)

type SearchPrivilegeFunc func(ctx context.Context, sessionID, privilegeCode string) error
type EndpointPrivilegesMiddlewareFunc func(f strictecho.StrictEchoHandlerFunc) strictecho.StrictEchoHandlerFunc

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
) strictecho.StrictEchoMiddlewareFunc {
	var endpointPrivileges = map[string]EndpointPrivilegesMiddlewareFunc{
		"Login":    without(),
		"GetUsers": withPrivilege(signingKey, searchPrivilegeFunc, "user_read"),
		"GetRoles": withPrivilege(signingKey, searchPrivilegeFunc, "role_read"),
	}

	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
		if op, ok := endpointPrivileges[operationID]; ok {
			return op(f)
		}

		return func(ctx echo.Context, request any) (any, error) {
			return nil, fmt.Errorf("not found")
		}
	}
}

// func abc(l *zap.Logger) strictecho.StrictEchoMiddlewareFunc {
// 	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
// 		return func(ctx echo.Context, request any) (any, error) {
// 			l.Info("StrictEchoMiddlewareFunc ABC",
// 				zap.String("operationID", operationID),
// 			)

// 			return f(ctx, request)
// 		}
// 	}
// }
