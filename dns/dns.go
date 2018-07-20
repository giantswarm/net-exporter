package dns

import (
	"fmt"
	"net"
	"sync"
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

	count      map[string]float64
	errorCount map[string]float64

	countMutex      sync.Mutex
	errorCountMutex sync.Mutex

	total      *prometheus.Desc
	errorTotal *prometheus.Desc
	latency    *prometheus.Desc
}

func NewCollector(config Config) (*DNSCollector, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.Logger must not be empty", config)
	}

	if len(config.Hosts) == 0 {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.Host must not be empty", config)
	}

	dnsCollector := &DNSCollector{
		logger: config.Logger,

		hosts: config.Hosts,

		count:      map[string]float64{},
		errorCount: map[string]float64{},

		total: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "check_total"),
			"Total number of DNS resolutions.",
			[]string{"host"},
			nil,
		),
		errorTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "check_error_total"),
			"Total number of DNS resolution errors.",
			[]string{"host"},
			nil,
		),
		latency: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "check_seconds"),
			"Time taken to resolve DNS.",
			[]string{"host"},
			nil,
		),
	}

	return dnsCollector, nil
}

func (c *DNSCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.errorTotal
	ch <- c.latency
}

func (c *DNSCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	for _, host := range c.hosts {
		wg.Add(1)

		go func(host string) {
			defer wg.Done()

			start := time.Now()

			c.countMutex.Lock()
			c.count[host]++
			c.countMutex.Unlock()

			_, err := net.LookupHost(host)
			if err != nil {
				c.logger.Log("level", "error", "message", "could not resolve dns", "host", host, "stack", fmt.Sprintf("%#v", err))

				c.errorCountMutex.Lock()
				c.errorCount[host]++
				ch <- prometheus.MustNewConstMetric(c.errorTotal, prometheus.CounterValue, c.errorCount[host], host)
				c.errorCountMutex.Unlock()

				return
			}

			elapsed := time.Since(start)

			c.countMutex.Lock()
			ch <- prometheus.MustNewConstMetric(c.total, prometheus.CounterValue, c.count[host], host)
			c.countMutex.Unlock()
			ch <- prometheus.MustNewConstMetric(c.latency, prometheus.GaugeValue, elapsed.Seconds(), host)
		}(host)
	}

	wg.Wait()
}
