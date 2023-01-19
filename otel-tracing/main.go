package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	otelSDKTrace "go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	rand.Seed(time.Now().UnixMicro())
	tracer := setupTracing()

	s := &service{
		tracer: tracer,
	}
	h := &handler{
		svc:    s,
		tracer: tracer,
	}
	http.HandleFunc("/hello", loggingMiddleware(tracingMiddleware(h.Handle, tracer), "manual-hello"))
	err := http.ListenAndServe("", nil)
	if err != nil {
		log.Println(err)
	}
}

func loggingMiddleware(next http.HandlerFunc, epName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s request completed in %s\n", epName, time.Since(start))
	}
}

func tracingMiddleware(next http.HandlerFunc, tracer trace.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "rest-request")
		defer span.End()
		r = r.WithContext(ctx)
		next(w, r)
	}
}

type handler struct {
	svc    *service
	tracer trace.Tracer
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "handler.Handle")
	defer span.End()
	_, err := w.Write([]byte(h.svc.hello(ctx)))
	if err != nil {
		panic(err)
	}
}

type service struct {
	tracer trace.Tracer
}

func (s *service) hello(ctx context.Context) string {
	_, span := s.tracer.Start(ctx, "service.hello")
	defer span.End()
	diceRoll := rand.Intn(100)
	fmt.Println("sleeping", diceRoll, "ms")
	span.SetAttributes(attribute.Int("sleep.duration.ms", int(diceRoll)))
	time.Sleep(time.Duration(diceRoll) * time.Millisecond)
	return "Hello World!"
}

func setupTracing() trace.Tracer {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		panic(err)
	}
	tp := otelSDKTrace.NewTracerProvider(
		// Always be sure to batch in production.
		otelSDKTrace.WithBatcher(exporter),
		// Record information about this application in a Resource.
		otelSDKTrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("hello-server"),
		)),
	)

	return tp.Tracer("http://github.com/t-margheim/demo-space/otel-tracing")
}
