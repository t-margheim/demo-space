package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/t-margheim/demo-space/metrics-poc/pkg/rsmetrics"
	"go.opentelemetry.io/otel/attribute"
	otlpExportPrometheus "go.opentelemetry.io/otel/exporters/prometheus"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	// provider, err := setupMetricsProvider()
	// if err != nil {
	// 	panic(err)
	// }
	// meter := provider.Meter("github.com/t-margheim/demo-space/metrics-poc")

	// err = rsmetrics.Initialize(meter)
	// if err != nil {
	// 	panic(err)
	// }

	h := &handler{}
	http.Handle("/hello", h)
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}

func setupMetricsProvider() (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpExportPrometheus.New()
	if err != nil {
		return nil, err
	}

	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	return provider, nil
}

type request struct {
	Name string
}

type response struct {
	Message string
}

type handler struct {
	svc service
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var ok bool
	start := time.Now()

	defer func() {
		okAttr := attribute.Bool("ok", ok)
		rsmetrics.Count(ctx, "requests_received", 1, okAttr)
		rsmetrics.Timing(ctx, "request_processed", time.Since(start), okAttr)
	}()

	var req request
	err := parseRequest(r.Body, &req)
	if err != nil {
		log.Println("failed to parseRequest:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := h.svc.sayHello(ctx, req.Name)

	err = writeResponse(w, response{Message: msg})
	if err != nil {
		log.Println("failed to writeResponse", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ok = true
}

type service struct{}

var rng = rand.New(rand.NewSource(time.Now().Unix()))

func (s *service) sayHello(ctx context.Context, name string) string {
	rsmetrics.Count(ctx,
		"a_in_names",
		int64(strings.Count(strings.ToLower(name), "a")),
		attribute.String("first_letter", string(name[0])),
	)
	rsmetrics.Count(ctx, "letters_in_names", int64(len(name)))

	time.Sleep(time.Duration(rng.Intn(1000)) * time.Millisecond)
	return fmt.Sprintf("Hello %s", name)
}
