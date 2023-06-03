package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func main() {
	_, cleanup, err := InitializeTracing("OCHaCafe sample")
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	tracer = otel.Tracer("OCHaCafe sample")

	ctx := context.Background()
	Operation1(ctx)
	Operation2(ctx)
	Operation3(ctx)
}

func Operation1(ctx context.Context) {
	_, span := tracer.Start(ctx, "op1")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
}

func Operation2(ctx context.Context) {
	_, span := tracer.Start(ctx, "op2")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
}

func Operation3(ctx context.Context) {
	_, span := tracer.Start(ctx, "op3")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
}

func InitializeTracing(serviceName string) (*sdktrace.TracerProvider, func(), error) {
	exporter, err := NewJaegerExporter()
	if err != nil {
		return nil, nil, err
	}

	r := NewResource(serviceName, "1.0.0", "local")
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := tp.ForceFlush(ctx); err != nil {
			log.Print(err)
		}
		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		if err := tp.Shutdown(ctx2); err != nil {
			log.Print(err)
		}
		cancel()
		cancel2()
	}
	return tp, cleanup, nil
}

func NewResource(serviceName string, version string, environment string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(version),
		attribute.String("environment", environment),
	)
}

func NewJaegerExporter() (sdktrace.SpanExporter, error) {
	endpoint := "localhost:14268/api/traces"

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func NewStdoutExporter() (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(os.Stderr),
	)
}
