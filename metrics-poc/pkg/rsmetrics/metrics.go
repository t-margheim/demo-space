package rsmetrics

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

var meter metric.Meter

var once = sync.Once{}

func setupMetrics(ctx context.Context) metric.Meter {
	// as long as no options are passed in, no error is possible
	exporter, _ := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
	)

	// labels/tags/resources that are common to all metrics.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
	)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			// collects and exports metric data every 30 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(30*time.Second)),
		),
	)

	return mp.Meter("")
}

func initialize() {
	log.Println("once.Do.initialize")
	if meter == nil {
		meter = setupMetrics(context.Background())
	}
	if proc == nil {
		proc = &standardProcessor{
			counters:   map[string]instrument.Int64Counter{},
			histograms: map[string]instrument.Int64Histogram{},
		}
	}
}

// Initialize is used to start the
func Initialize(m metric.Meter) error {
	if m == nil {
		return errors.New("invalid value passed as metric.Meter: nil")
	}

	once.Do(func() {
		log.Println("Initialize.once.Do")
		meter = m
		proc = &standardProcessor{
			counters:   map[string]instrument.Int64Counter{},
			histograms: map[string]instrument.Int64Histogram{},
		}
	})

	return nil
}

func Count(ctx context.Context, metricName string, value int64, attrs ...attribute.KeyValue) {
	once.Do(initialize)
	proc.count(ctx, metricName, value, attrs...)
}

func Timing(ctx context.Context, metricName string, duration time.Duration, attrs ...attribute.KeyValue) {
	once.Do(initialize)
	proc.timing(ctx, metricName, duration, attrs...)
}
