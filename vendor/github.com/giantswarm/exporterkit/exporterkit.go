package exporterkit

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

// Config if the configuration to create an Exporter.
type Config struct {
	Collectors []prometheus.Collector
	Logger     micrologger.Logger
}

// Exporter runs a slice of Prometheus Collectors.
type Exporter struct {
	collectors []prometheus.Collector
	logger     micrologger.Logger
}

// New creates a new Exporter, given a Config.
func New(config Config) (*Exporter, error) {
	if config.Collectors == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Collectors must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	exporter := Exporter{
		collectors: config.Collectors,
		logger:     config.Logger,
	}

	return &exporter, nil
}

// Run starts the Exporter.
func (e *Exporter) Run() {
	prometheus.MustRegister(e.collectors...)

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok\n")
	}))

	http.ListenAndServe("0.0.0.0:8000", nil)
}
