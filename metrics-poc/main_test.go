package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t-margheim/demo-space/metrics-poc/pkg/rsmetrics"
	"go.opentelemetry.io/otel/attribute"
)

func Test_service_sayHello(t *testing.T) {
	tests := []struct {
		name            string
		s               *service
		want            string
		wantMockMetrics *rsmetrics.MockProcessor
	}{
		{
			name: "Timothy",
			s:    &service{},
			want: "Hello Timothy",
			wantMockMetrics: &rsmetrics.MockProcessor{
				CountCalledTimes: 2,
				CountCalledWithNames: []string{
					"a_in_names",
					"letters_in_names",
				},
				CountCalledWithValues: []int64{
					0,
					7,
				},
				CountCalledWithAttributes: [][]attribute.KeyValue{
					{attribute.String("first_letter", "T")},
					nil,
				},
			},
		},
		{
			name: "Madagascar",
			s:    &service{},
			want: "Hello Madagascar",
			wantMockMetrics: &rsmetrics.MockProcessor{
				CountCalledTimes: 2,
				CountCalledWithNames: []string{
					"a_in_names",
					"letters_in_names",
				},
				CountCalledWithValues: []int64{
					4,
					10,
				},
				CountCalledWithAttributes: [][]attribute.KeyValue{
					{attribute.String("first_letter", "M")},
					nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMetrics := rsmetrics.Mock()
			s := &service{}
			if got := s.sayHello(context.Background(), tt.name); got != tt.want {
				t.Errorf("service.sayHello() = %v, want %v", got, tt.want)
			}
			assert.Equal(t, tt.wantMockMetrics, mockMetrics, "unexpected metrics data")
		})
	}
}
