package httptransport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
	"go.uber.org/zap"
)

type SearchPrivilegeFunc func(ctx context.Context, sessionID, privilegeCode string) error
type EndpointPrivilegesMiddlewareFunc func(f strictecho.StrictEchoHandlerFunc) strictecho.StrictEchoHandlerFunc
type SessionService interface {
	Search(ctx context.Context, sessionID, privilegeCode string) error
}

func endpointPrivilegeKey(c echo.Context) string {
	return fmt.Sprintf("%s%s", strings.ToLower(c.Request().Method), c.Path())
}

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
	sessionService SessionService,
	signingKey string,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			endpointPrivilegeKey := endpointPrivilegeKey(c)

			if _, ok := endpointWithPrivileges[endpointPrivilegeKey]; !ok {
				if _, ok = endpointWithout[endpointPrivilegeKey]; !ok {
					return fmt.Errorf("privilege not found")
				}
			}

			return next(c)
		}
	}
}

// func AuthorizationMiddleware(
// 	signingKey string,
// 	sessionService SessionService,
// ) strictecho.StrictEchoMiddlewareFunc {
// 	var endpointPrivileges = EndpointPrivilegesMiddlewareFuncs(signingKey, sessionService)

// 	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
// 		if op, ok := endpointPrivileges[operationID]; ok {
// 			return op(f)
// 		}

// 		return func(ctx echo.Context, request any) (any, error) {
// 			return nil, fmt.Errorf("not found")
// 		}
// 	}
// }

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
