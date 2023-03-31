package telemetry

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"github.com/rs/zerolog/log"
)

type options struct {
	serviceName     string
	version         string
	env             string
	metricsEndpoint string
	tracesEndpoint  string
	ctx             context.Context
	configure       bool
}

func (o *options) MetricsEnabled() bool {
	return o.metricsEndpoint != ""
}

func (o *options) TracesEnabled() bool {
	return o.tracesEndpoint != ""
}

type Option func(o *options) *options

var (
	telemetryOpts *options
)

func InitTelemetry(opts ...Option) {
	opt := &options{
		serviceName:     "default-name",
		version:         "v0.0.0",
		env:             "dev",
		metricsEndpoint: "",
		tracesEndpoint:  "",
		ctx:             context.Background(),
	}
	for _, f := range opts {
		opt = f(opt)
	}
	log.Debug().Msg("starting telemetry providers")

	telemetryOpts = opt

	initMetrics(opt)
	initTracer(opt)

	log.Debug().Msg("ending telemetry providers")
}

// WithServiceName configures service name
func WithServiceName(n string) Option {
	return func(o *options) *options {
		o.serviceName = n
		return o
	}
}

// WithVersion configures service version
func WithVersion(v string) Option {
	return func(o *options) *options {
		o.version = v
		return o
	}
}

// WithEnvironment configures environment name
func WithEnvironment(e string) Option {
	return func(o *options) *options {
		o.env = e
		return o
	}
}

// WithMetricsEndpoint configures otel metrics endpoint
func WithMetricsEndpoint(e string) Option {
	return func(o *options) *options {
		o.metricsEndpoint = e
		return o
	}
}

// WithTracesEndpoint configures otel traces endpoint
func WithTracesEndpoint(e string) Option {
	return func(o *options) *options {
		o.tracesEndpoint = e
		return o
	}
}

// WithContext set's up the configuration context
func WithContext(c context.Context) Option {
	return func(o *options) *options {
		o.ctx = c
		return o
	}
}

func SetupRestTracing(r *chi.Mux) {
	if telemetryOpts == nil {
		log.Print("telemetry configuration not started, please call telemetry.InitTelemetry before instrument your code")
		return
	}
	r.Use(otelchi.Middleware(telemetryOpts.serviceName, otelchi.WithChiRoutes(r)))
}
