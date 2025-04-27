package tracing

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTraceRatio = 1
)

type TracerShutdown func(context.Context) error

type TracingClientConfig struct {
	// OTLP Exporter; example: 10.1.30.53:4317
	Endpoint string
	// example: auth-v2
	ServiceName string
	// example: dev
	Environment string
}

type TracingClientOption func(*tracingClientConfig)

type tracingClientConfig struct {
	// default 1.0
	TraceIDRatio float64

	// default insecure.NewCredentials()
	Cred credentials.TransportCredentials
}

func WithTracingRatio(ratio float64) TracingClientOption {
	return func(config *tracingClientConfig) {
		config.TraceIDRatio = ratio
	}
}

// func ProvideTracer(ctx context.Context, cfg *conf.Config) (TracerShutdown, error) {
// 	sd, err := initOpentelemetryTracing(ctx,
// 		TracingClientConfig{
// 			Endpoint:    cfg.Tracing.Address,
// 			ServiceName: cfg.ServiceName,
// 			Environment: cfg.Environment,
// 		},
// 		WithTracingRatio(cfg.Tracing.Ratio),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while init Opentelemetry tracing | %w", err)
// 	}

// 	return sd, nil
// }

// GetTraceID извлекает TraceID из контекста.
func GetTraceID(ctx context.Context) string {
	if trace.SpanContextFromContext(ctx).IsValid() {
		return trace.SpanContextFromContext(ctx).TraceID().String()
	}

	return ""
}

func initOpentelemetryTracing(ctx context.Context, cfg TracingClientConfig, opts ...TracingClientOption) (func(context.Context) error, error) {
	optCfg := tracingClientConfig{
		TraceIDRatio: defaultTraceRatio,
		Cred:         insecure.NewCredentials(),
	}

	for _, opt := range opts {
		opt(&optCfg)
	}

	if cfg.Endpoint == "" || cfg.Environment == "" || cfg.ServiceName == "" {
		return nil, errors.New("opentelemetry error: values: Endpoint, Environment, ServiceName can't be empty string")
	}

	conn, err := grpc.NewClient(cfg.Endpoint, grpc.WithTransportCredentials(optCfg.Cred))
	if err != nil {
		return nil, fmt.Errorf("grpc new client error | %w", err)
	}

	client := otlptracegrpc.NewClient(otlptracegrpc.WithGRPCConn(conn))

	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			attribute.String("environment", cfg.Environment),
			// attribute.String("zone", getZone()),
			// attribute.String("k8s_host", getK8sComputeHost()),
		)),
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(optCfg.TraceIDRatio))),
	)

	otel.SetTextMapPropagator(b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)))

	otel.SetTracerProvider(tp)

	return func(ctx context.Context) error {
		var errs error

		if err := tp.ForceFlush(ctx); err != nil {
			errs = errors.Join(err)
		}

		if err := tp.Shutdown(ctx); err != nil {
			errs = errors.Join(err)
		}

		if err := conn.Close(); err != nil {
			errs = errors.Join(err)
		}

		return errs
	}, nil
}
