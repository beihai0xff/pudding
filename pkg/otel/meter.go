// Package otel provides OpenTelemetry utilities.
package otel

import (
	"sync"

	prometheusExporter "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"

	"github.com/beihai0xff/pudding/pkg/log"
)

var once sync.Once
var meterProvider metric.MeterProvider

// GetMeterProvider get the global meter provider
func GetMeterProvider() metric.MeterProvider {
	once.Do(func() {
		export, err := prometheusExporter.New()
		if err != nil {
			log.Panicf("failed to create prometheus exporter: %v", err)
		}
		meterProvider = metricsdk.NewMeterProvider(metricsdk.WithReader(export.Reader))
	})

	return meterProvider
}
