package sessionsvc

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter           = otel.Meter("sessionsvc") //nolint:gochecknoglobals
	authCallCounter metric.Int64Counter        //nolint:gochecknoglobals
)

func init() { //nolint:gochecknoinits
	totalCounter, err := meter.Int64Counter(
		"auth_calls",
		metric.WithDescription("The count of auth calls"),
		metric.WithUnit(""),
	)
	if err != nil {
		panic(fmt.Errorf("error while create auth_calls metric | %w", err))
	}

	authCallCounter = totalCounter
}

// Инкремент счетчика успешного логина
func incrLoginSuccess(ctx context.Context) {
	authCallCounter.Add(
		ctx,
		1,
		metric.WithAttributeSet(
			attribute.NewSet(
				attribute.Bool("success", true),
			),
		),
	)
}

// Инкремент счетчика неудачного логина
func incrLoginFail(ctx context.Context, callKind string) {
	authCallCounter.Add(
		ctx,
		1,
		metric.WithAttributeSet(
			attribute.NewSet(
				attribute.Bool("success", false),
				attribute.String("call_kind", strings.ToLower(callKind)),
			),
		),
	)
}
