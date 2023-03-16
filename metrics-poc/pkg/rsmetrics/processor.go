package rsmetrics

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
)

type processor interface {
	count(ctx context.Context, metricName string, value int64, attrs ...attribute.KeyValue)
	timing(ctx context.Context, metricName string, duration time.Duration, attrs ...attribute.KeyValue)
}

var proc processor

type standardProcessor struct {
	counters   map[string]instrument.Int64Counter
	histograms map[string]instrument.Int64Histogram

	countersLock   sync.Mutex
	histogramsLock sync.Mutex
}

func (p *standardProcessor) count(ctx context.Context, metricName string, value int64, attrs ...attribute.KeyValue) {
	p.countersLock.Lock()
	defer p.countersLock.Unlock()
	counter, ok := p.counters[metricName]
	if !ok {
		// can't get error here?
		counter, _ = meter.Int64Counter(metricName)
		p.counters[metricName] = counter
	}
	counter.Add(ctx, value, attrs...)
}

func (p *standardProcessor) timing(ctx context.Context, metricName string, duration time.Duration, attrs ...attribute.KeyValue) {
	p.histogramsLock.Lock()
	defer p.histogramsLock.Unlock()
	histogram, ok := p.histograms[metricName]
	if !ok {
		// can't get error here?
		histogram, _ = meter.Int64Histogram(metricName, instrument.WithUnit("milliseconds"))
		p.histograms[metricName] = histogram
	}
	histogram.Record(ctx, int64(duration.Milliseconds()), attrs...)
}
