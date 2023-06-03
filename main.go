package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func main() {
	_, cleanup, err := NewTracerProvider("OCHaCafe sample")
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
