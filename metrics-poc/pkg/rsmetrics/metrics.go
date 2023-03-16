package rsmetrics

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

var meter metric.Meter

// Initialize is used to start the
func Initialize(m metric.Meter) error {
	if m == nil {
		return errors.New("invalid value passed as metric.Meter: nil")
	}
	meter = m
	proc = &standardProcessor{
		counters:   map[string]instrument.Int64Counter{},
		histograms: map[string]instrument.Int64Histogram{},
	}
	return nil
}

func Count(ctx context.Context, metricName string, value int64, attrs ...attribute.KeyValue) {
	proc.count(ctx, metricName, value, attrs...)
}

func Timing(ctx context.Context, metricName string, duration time.Duration, attrs ...attribute.KeyValue) {
	proc.timing(ctx, metricName, duration, attrs...)
}
