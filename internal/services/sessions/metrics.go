package sessionsvc

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func (s *SessionSvc) incrLoginFail(ctx context.Context, failure_kind string) {
	s.metricsLoginCounter.Add(
		ctx,
		1,
		metric.WithAttributeSet(
			attribute.NewSet(
				attribute.Bool("success", false),
				attribute.String("failure_kind", strings.ToLower(failure_kind)),
			),
		),
	)
}

func (s *SessionSvc) incrLoginSuccess(ctx context.Context) {
	s.metricsLoginCounter.Add(
		ctx,
		1,
		metric.WithAttributeSet(
			attribute.NewSet(
				attribute.Bool("success", true),
			),
		),
	)
}
