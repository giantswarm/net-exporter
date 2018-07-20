package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/giantswarm/exporterkit"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "dns"
)

type Config struct {
	Logger micrologger.Logger

	Hosts []string
}

type DNSCollector struct {
	logger micrologger.Logger

	hosts []string

	// Internals for calculating error count.
	errorCountDesc *prometheus.Desc
	errorCount     map[string]float64

	// Internals for calculating histograms.
	latencyHistogramsDesc *prometheus.Desc
	latencyCount          map[string]uint64
	latencySum            map[string]float64
	latencyBuckets        map[string]map[float64]uint64
}

func NewCollector(config Config) (*DNSCollector, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.Logger must not be empty", config)
	}

	if len(config.Hosts) == 0 {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.Host must not be empty", config)
	}

	latencyBuckets := map[string]map[float64]uint64{}
	for _, host := range config.Hosts {
		latencyBuckets[host] = map[float64]uint64{
			0.01: 0,
			0.1:  0,
			0.5:  0,
			1:    0,
			5:    0,
			10:   0,
		}
	}

	dnsCollector := &DNSCollector{
		logger: config.Logger,

		hosts: config.Hosts,

		errorCountDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "resolution_error_count"),
			"Total number of DNS resolution errors.",
			[]string{"host"},
			nil,
		),
		errorCount: map[string]float64{},

		latencyHistogramsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "resolution_seconds"),
			"A histogram of the DNS resolution durations.",
			[]string{"host"},
			nil,
		),
		latencyCount:   map[string]uint64{},
		latencySum:     map[string]float64{},
		latencyBuckets: latencyBuckets,
	}

	return dnsCollector, nil
}

func (c *DNSCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.errorCountDesc
	ch <- c.latencyHistogramsDesc
}

func (c *DNSCollector) Collect(ch chan<- prometheus.Metric) {

	for _, host := range c.hosts {
		start := time.Now()

		_, err := net.LookupHost(host)
		if err != nil {
			c.logger.Log("level", "error", "message", "could not resolve dns", "host", host, "stack", fmt.Sprintf("%#v", err))
			c.errorCount[host]++
			ch <- prometheus.MustNewConstMetric(c.errorCountDesc, prometheus.CounterValue, c.errorCount[host], host)
		}

		elapsed := float64(time.Since(start).Seconds())

		c.latencyCount[host]++
		c.latencySum[host] += elapsed

		for bucket, _ := range c.latencyBuckets[host] {
			if elapsed <= bucket {
				c.latencyBuckets[host][bucket]++
			}
		}

		ch <- prometheus.MustNewConstHistogram(
			c.latencyHistogramsDesc,
			c.latencyCount[host],
			c.latencySum[host],
			c.latencyBuckets[host],
			host,
		)
	}
}
