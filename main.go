package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("TodoAPI"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "http://localhost:14268/api/traces",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

var tracer trace.Tracer

func main() {
	ctx := context.Background()

	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)

	// router 設定
	r := gin.Default()
	r.Use(otelgin.Middleware("OCHaCafe"))
	r.GET("/ochacafe", ochacafe1)

	fmt.Println("server start...")
	r.Run(":3000")
}

func ochacafe1(c *gin.Context) {
	ctx := c.Request.Context()
	_, span := tracer.Start(ctx, "ochacafe1")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
	ochacafe2(ctx)
}

func ochacafe2(ctx context.Context) {
	_, span := tracer.Start(ctx, "ochacafe2")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
	ochacafe3(ctx)
}

func ochacafe3(ctx context.Context) {
	_, span := tracer.Start(ctx, "ochacafe3")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
}
