package sessionsvc

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func (s *SessionSvc) incrLoginFail(ctx context.Context, kind string) {
	s.metricsLoginCounter.Add(
		ctx,
		1,
		metric.WithAttributeSet(
			attribute.NewSet(
				attribute.Bool("success", false),
				attribute.String("kind", strings.ToLower(kind)),
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
