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
)

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	// resource の生成
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("OCHaCafe Service"),
			semconv.ServiceInstanceIDKey.String("Local"),
			semconv.ServiceVersionKey.String("1.0.0-test"),
			semconv.TelemetrySDKNameKey.String("OpenTelemetry~~~!!!!"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// traceExporter の生成
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// spanProcessor の生成
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	// traceProvider の生成
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// tracer の設定
	otel.SetTracerProvider(tracerProvider)

	// propagater の生成
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

var tracer = otel.Tracer("OCHaCafe Trace Sample APP PKG")

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
	r.Run(":8080")
}

func ochacafe1(c *gin.Context) {
	ctx := c.Request.Context()
	_, span := tracer.Start(ctx, "ochacafe span 1")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
	ochacafe2(c)
}

func ochacafe2(c *gin.Context) {
	ctx := c.Request.Context()
	_, span := tracer.Start(ctx, "ochacafe span 2")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
	ochacafe3(c)
}

func ochacafe3(c *gin.Context) {
	ctx := c.Request.Context()
	_, span := tracer.Start(ctx, "ochacafe span 3")
	defer span.End()
	time.Sleep(100 * time.Millisecond)
}
