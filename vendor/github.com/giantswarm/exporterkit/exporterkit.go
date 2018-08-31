package exporterkit

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/giantswarm/microendpoint/endpoint/healthz"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// DefaultAddress is the address the Exporter will run on if no address is configured.
	DefaultAddress = "http://0.0.0.0:8000"
)

// Config if the configuration to create an Exporter.
type Config struct {
	Address    string
	Collectors []prometheus.Collector
	Logger     micrologger.Logger
}

// Exporter runs a slice of Prometheus Collectors.
type Exporter struct {
	address    string
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

	if config.Address == "" {
		config.Address = DefaultAddress
	}

	exporter := Exporter{
		address:    config.Address,
		collectors: config.Collectors,
		logger:     config.Logger,
	}

	return &exporter, nil
}

// Run starts the Exporter.
func (e *Exporter) Run() {
	var err error

	var healthzEndpoint *healthz.Endpoint
	{
		c := healthz.Config{
			Logger: e.logger,
		}
		healthzEndpoint, err = healthz.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var newServer server.Server
	{
		c := server.Config{
			Logger:        e.logger,
			Endpoints:     []server.Endpoint{healthzEndpoint},
			ListenAddress: e.address,
		}

		newServer, err = server.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	prometheus.MustRegister(e.collectors...)

	go newServer.Boot()

	listener := make(chan os.Signal, 2)
	signal.Notify(listener, os.Interrupt, os.Kill)

	<-listener

	go func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			newServer.Shutdown()
		}()

		os.Exit(0)
	}()

	<-listener

	os.Exit(0)
}
