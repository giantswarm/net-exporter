package dns

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/giantswarm/exporterkit/histogramvec"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "dns"

	bucketStart  = 0.001
	bucketFactor = 2
	numBuckets   = 15
)

type Config struct {
	Logger micrologger.Logger

	Hosts []string
}

type DNSCollector struct {
	logger micrologger.Logger

	hosts []string

	latencyHistogramVec  *histogramvec.HistogramVec
	latencyHistogramDesc *prometheus.Desc

	errorTotal     float64
	errorTotalDesc *prometheus.Desc
}

func NewCollector(config Config) (*DNSCollector, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if len(config.Hosts) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.Host must not be empty", config)
	}

	var err error
	var latencyHistogramVec *histogramvec.HistogramVec
	{
		c := histogramvec.Config{
			BucketLimits: prometheus.ExponentialBuckets(bucketStart, bucketFactor, numBuckets),
		}
		latencyHistogramVec, err = histogramvec.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	dnsCollector := &DNSCollector{
		logger: config.Logger,

		hosts: config.Hosts,

		latencyHistogramVec: latencyHistogramVec,
		latencyHistogramDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "latency_seconds"),
			"Histogram of latency of DNS resolutions.",
			[]string{"host"},
			nil,
		),

		errorTotalDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "error_total"),
			"Total of DNS resolution errors.",
			nil,
			nil,
		),
	}

	return dnsCollector, nil
}

func (c *DNSCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.latencyHistogramDesc
	ch <- c.errorTotalDesc
}

func (c *DNSCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	for _, host := range c.hosts {
		wg.Add(1)

		go func(host string) {
			defer wg.Done()

			start := time.Now()

			_, err := net.LookupHost(host)
			if err != nil {
				c.logger.Log("level", "error", "message", "could not resolve dns", "host", host, "stack", fmt.Sprintf("%#v", err))
				c.errorTotal++
				return
			}

			elapsed := time.Since(start)

			c.latencyHistogramVec.Add(host, elapsed.Seconds())
		}(host)
	}

	wg.Wait()

	c.latencyHistogramVec.Ensure(c.hosts)

	ch <- prometheus.MustNewConstMetric(c.errorTotalDesc, prometheus.CounterValue, c.errorTotal)
	for host, histogram := range c.latencyHistogramVec.Histograms() {
		ch <- prometheus.MustNewConstHistogram(
			c.latencyHistogramDesc,
			histogram.Count(), histogram.Sum(), histogram.Buckets(),
			host,
		)
	}
}
