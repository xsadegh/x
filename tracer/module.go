package tracer

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
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
	Logs     string `yaml:"logs"`
	Version  string `yaml:"version"`
	Service  string `yaml:"service"`
	Endpoint string `yaml:"endpoint"`
	Protocol string `yaml:"protocol"`
	Insecure bool   `yaml:"insecure"`

	Headers map[string]string `yaml:"headers"`
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
	var exporter *otlptrace.Exporter
	if config.Protocol == "grpc" {
		var opts []otlptracegrpc.Option
		if config.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		opts = append(opts, otlptracegrpc.WithTimeout(5*time.Second))
		opts = append(opts, otlptracegrpc.WithHeaders(config.Headers))
		opts = append(opts, otlptracegrpc.WithEndpoint(config.Endpoint))
		opts = append(opts, otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:        true,
			MaxInterval:    2 * time.Second,
			MaxElapsedTime: 10 * time.Second,
		}))
		exporter, _ = otlptrace.New(
			context.Background(),
			otlptracegrpc.NewClient(opts...),
		)
	} else {
		var opts []otlptracehttp.Option
		if config.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		opts = append(opts, otlptracehttp.WithTimeout(5*time.Second))
		opts = append(opts, otlptracehttp.WithHeaders(config.Headers))
		opts = append(opts, otlptracehttp.WithEndpoint(config.Endpoint))
		opts = append(opts, otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:        true,
			MaxInterval:    2 * time.Second,
			MaxElapsedTime: 10 * time.Second,
		}))

		exporter, _ = otlptrace.New(
			context.Background(),
			otlptracehttp.NewClient(opts...),
		)
	}
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
