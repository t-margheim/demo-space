package rsmetrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// Mock would be intended for unit test spy assertions.
func Mock() *MockProcessor {
	p := &MockProcessor{}
	proc = p
	return p
}

type MockProcessor struct {
	CountCalledTimes          int
	CountCalledWithNames      []string
	CountCalledWithValues     []int64
	CountCalledWithAttributes [][]attribute.KeyValue

	TimingCalledTimes          int
	TimingCalledWithNames      []string
	TimingCalledWithDurations  []time.Duration
	TimingCalledWithAttributes [][]attribute.KeyValue
}

func (p *MockProcessor) count(ctx context.Context, metricName string, value int64, attrs ...attribute.KeyValue) {
	p.CountCalledTimes++
	p.CountCalledWithNames = append(p.CountCalledWithNames, metricName)
	p.CountCalledWithValues = append(p.CountCalledWithValues, value)
	p.CountCalledWithAttributes = append(p.CountCalledWithAttributes, attrs)
}

func (p *MockProcessor) timing(ctx context.Context, metricName string, duration time.Duration, attrs ...attribute.KeyValue) {
	p.TimingCalledTimes++
	p.TimingCalledWithNames = append(p.TimingCalledWithNames, metricName)
	p.TimingCalledWithDurations = append(p.TimingCalledWithDurations, duration)
	p.TimingCalledWithAttributes = append(p.TimingCalledWithAttributes, attrs)
}
