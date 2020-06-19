package ntp

import (
	"fmt"
	"sync"
	"time"

	"github.com/beevik/ntp"
	"github.com/giantswarm/exporterkit/histogramvec"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
)

const (
	namespace = "ntp"

	bucketStart  = 0.001
	bucketFactor = 2
	numBuckets   = 15
)

// Config provides the necessary configuration for creating a Collector.
type Config struct {
	Logger micrologger.Logger

	NTPServers []string
	SourceHost string
}

// Collector implements the Collector interface, exposing DNS latency information.
type Collector struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	ntpServers []string
	sourceHost string

	latencyHistogramVec  *histogramvec.HistogramVec
	latencyHistogramDesc *prometheus.Desc

	errorCount     prometheus.Counter
	syncErrorCount *prometheus.CounterVec
}

// New creates a Collector, given a Config.
func New(config Config) (*Collector, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if len(config.NTPServers) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.NTPServers must not be empty", config)
	}

	if len(config.SourceHost) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.SourceHost must not be empty", config)
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

	errorCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "error_total"),
		Help: "Total number of internal errors.",
	})
	syncErrorCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, "", "ntp_sync_error_total"),
			Help: "Total number of errors ntp syncs.",
		},
		[]string{"ntp-server"},
	)

	prometheus.MustRegister(errorCount)
	prometheus.MustRegister(syncErrorCount)

	collector := &Collector{
		logger: config.Logger,

		ntpServers: config.NTPServers,
		sourceHost: config.SourceHost,

		latencyHistogramVec: latencyHistogramVec,
		latencyHistogramDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "latency_seconds"),
			"Histogram of latency of NTP sync requests.",
			[]string{"ntp-servers"},
			nil,
		),

		errorCount:     errorCount,
		syncErrorCount: syncErrorCount,
	}

	return collector, nil
}

// Describe implements the Describe method of the Collector interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.latencyHistogramDesc
}

func (c *Collector) ntpsync(host string, ntpServer string, latencyHistogramVec *histogramvec.HistogramVec) {
	start := time.Now()

	_, err := ntp.Time(ntpServer)
	if err != nil {
		c.logger.Log("level", "error", "message", fmt.Sprintf("could not sync time with ntp server %#q from host %#q", ntpServer, host), "stack", microerror.JSON(err))
		c.syncErrorCount.WithLabelValues(host).Inc()
		return
	}

	elapsed := time.Since(start)

	err = latencyHistogramVec.Add(ntpServer, elapsed.Seconds())
	if err != nil {
		c.logger.Log("level", "error", "message", fmt.Sprintf("failed to update latency histogram for host %#q", host), "stack", microerror.JSON(err))
		c.syncErrorCount.WithLabelValues(host, ntpServer).Inc()
		return
	}
}

// Collect implements the Collect method of the Collector interface.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	for _, ntpServer := range c.ntpServers {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()

			c.ntpsync(c.sourceHost, ntpServer, c.latencyHistogramVec)
		}(ntpServer)
	}

	wg.Wait()

	c.latencyHistogramVec.Ensure(c.ntpServers)

	for ntpServer, histogram := range c.latencyHistogramVec.Histograms() {
		ch <- prometheus.MustNewConstHistogram(
			c.latencyHistogramDesc,
			histogram.Count(), histogram.Sum(), histogram.Buckets(),
			ntpServer,
		)
	}
}
