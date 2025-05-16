package tracer

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

var MODULE = fx.Module(
	"TRACER",
	fx.Provide(NewTracer),
)

type Config struct {
	Version  string `yaml:"version"`
	Service  string `yaml:"service"`
	Endpoint string `yaml:"endpoint"`
}

type Tracer interface {
	Stop() error
	Tracer() trace.Tracer
}

type tracer struct {
	tracer   trace.Tracer
	exporter *otlptrace.Exporter
}

func NewTracer(config Config) Tracer {
	exporter, _ := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(config.Endpoint),
			otlptracegrpc.WithTimeout(5*time.Second),
			otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
				Enabled:        true,
				MaxInterval:    2 * time.Second,
				MaxElapsedTime: 10 * time.Second,
			}),
		),
	)
	resources, _ := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.Service),
			semconv.ServiceVersionKey.String(config.Version),
		),
	)

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)

	otel.SetTracerProvider(provider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)

	return &tracer{exporter: exporter, tracer: provider.Tracer(config.Service)}
}

func (t *tracer) Stop() error {
	return t.exporter.Shutdown(context.Background())
}

func (t *tracer) Tracer() trace.Tracer {
	return t.tracer
}
