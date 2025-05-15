package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

var MODULE = fx.Module(
	"TRACER",
	fx.Provide(NewTracer),
)

type Config struct {
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
		),
	)
	resources, _ := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", config.Service),
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
