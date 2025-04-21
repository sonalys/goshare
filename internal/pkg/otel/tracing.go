package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const name = "github.com/sonalys/goshare"

var (
	Tracer = otel.Tracer(name)
	Meter  = otel.Meter(name)
	Logger = global.Logger(name)
)

type Provider struct{}

func (tp *Provider) TracerProvider() trace.TracerProvider {
	return otel.GetTracerProvider()
}

func (tp *Provider) MeterProvider() metric.MeterProvider {
	return otel.GetMeterProvider()
}

func (tp *Provider) TextMapPropagator() propagation.TextMapPropagator {
	return otel.GetTextMapPropagator()
}
