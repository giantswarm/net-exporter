package network

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/giantswarm/exporterkit/histogramvec"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	namespace = "network"

	bucketStart  = 0.001
	bucketFactor = 2
	numBuckets   = 5

	// numNeighbours is the number of neighbours for the net-exporter to dial.
	// The lower the number, the higher the likelihood that a net-exporter is not dialed
	// in case of failures - e.g: if numNeighbours is 1, if a single net-exporter is down,
	// its neighbour will not be pinged.
	// The higher the number, the higher the cardinality of network latency metrics exposed
	// by the net-exporter.
	// Having a value of 2 means that 2 specific net-exporters need to be down
	// for one net-exporter to not be dialed, without exposing very high cardinality metrics.
	numNeighbours = 2
)

// Config provides the necessary configuration for creating a Collector.
type Config struct {
	Dialer    *net.Dialer
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	Namespace string
	Port      string
	Service   string
}

// Collector implements the Collector interface, exposing network latency information.
type Collector struct {
	dialer    *net.Dialer
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	namespace string
	port      string
	service   string

	// scrapeID is used to identify logs for a Collect call.
	scrapeID uint64

	latencyHistogramVec  *histogramvec.HistogramVec
	latencyHistogramDesc *prometheus.Desc

	errorCount     prometheus.Counter
	dialErrorCount *prometheus.CounterVec
}

// New creates a Collector, given a Config.
func New(config Config) (*Collector, error) {
	if config.Dialer == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Dialer must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.Namespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Namespace must not be empty", config)
	}
	if config.Port == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Port must not be empty", config)
	}
	if config.Service == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Service must not be empty", config)
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
	dialErrorCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, "", "dial_error_total"),
			Help: "Total number of errors dialing hosts.",
		},
		[]string{"host"},
	)
	prometheus.MustRegister(errorCount)
	prometheus.MustRegister(dialErrorCount)

	collector := &Collector{
		dialer:    config.Dialer,
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		namespace: config.Namespace,
		port:      config.Port,
		service:   config.Service,

		latencyHistogramVec: latencyHistogramVec,
		latencyHistogramDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "latency_seconds"),
			"Histogram of latency of network dials.",
			[]string{"host"},
			nil,
		),

		errorCount:     errorCount,
		dialErrorCount: dialErrorCount,
	}

	return collector, nil
}

// Describe implements the Describe method of the Collector interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.latencyHistogramDesc
}

// Collect implements the Collect method of the Collector interface.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	atomic.AddUint64(&c.scrapeID, 1)

	scrapingStart := time.Now()
	c.logger.Log("level", "info", "message", "collecting metrics", "scrapeID", c.scrapeID)

	c.logger.Log("level", "info", "message", "collecting service from kubernetes api", "service", c.service, "scrapeID", c.scrapeID)
	service, err := c.k8sClient.CoreV1().Services(c.namespace).Get(c.service, metav1.GetOptions{})
	if err != nil {
		c.logger.Log("level", "error", "message", "could not collect service from kubernetes api", "scrapeID", c.scrapeID, "stack", microerror.Stack(err))
		c.errorCount.Inc()
		return
	}

	c.logger.Log("level", "info", "message", "collected service from kubernetes api", "service ", c.service, "scrapeID", c.scrapeID)

	c.logger.Log("level", "info", "message", "collecting endpoints for service from kubernetes api", "service", c.service, "scrapeID", c.scrapeID)
	endpoints, err := c.k8sClient.CoreV1().Endpoints(c.namespace).Get(c.service, metav1.GetOptions{})
	if err != nil {
		c.logger.Log("level", "error", "message", "could not collect endpoints for service from kubernetes api ", "service", c.service, "scrapeID", c.scrapeID, "stack", microerror.Stack(err))
		c.errorCount.Inc()
		return
	}

	c.logger.Log("level", "info", "message", "collected endpoints for service from kubernetes api", "service", c.service, "scrapeID", c.scrapeID)

	hosts := []string{}
	hosts = append(hosts, fmt.Sprintf("%v:%v", service.Spec.ClusterIP, c.port))

	neighbours, err := c.getNeighbours(numNeighbours, endpoints.Subsets)
	if err != nil {
		c.logger.Log("level", "error", "message", "could not get neighbours", "service", c.service, "scrapeID", c.scrapeID, "stack", microerror.Stack(err))
		c.errorCount.Inc()
		return
	}
	for _, neighbour := range neighbours {
		hosts = append(hosts, fmt.Sprintf("%v:%v", neighbour, c.port))
	}

	var wg sync.WaitGroup

	for _, host := range hosts {
		wg.Add(1)

		go func(host string) {
			defer wg.Done()

			start := time.Now()

			c.logger.Log("level", "info", "message", "dialing host", "host", host, "scrapeID", c.scrapeID)
			conn, err := c.dialer.Dial("tcp", host)
			if err != nil {
				c.logger.Log("level", "error", "message", "could not dial host", "host", host, "scrapeID", c.scrapeID, "stack", microerror.Stack(err))
				c.dialErrorCount.WithLabelValues(host).Inc()
				return
			}
			defer conn.Close()

			elapsed := time.Since(start)
			c.logger.Log("level", "info", "message", "dialed host", "host", host, "scrapeTime", elapsed.Seconds(), "scrapeID", c.scrapeID)

			c.latencyHistogramVec.Add(host, elapsed.Seconds())
		}(host)
	}

	wg.Wait()

	c.latencyHistogramVec.Ensure(hosts)

	for host, histogram := range c.latencyHistogramVec.Histograms() {
		ch <- prometheus.MustNewConstHistogram(
			c.latencyHistogramDesc,
			histogram.Count(), histogram.Sum(), histogram.Buckets(),
			host,
		)
	}

	scrapingElapsed := time.Since(scrapingStart)
	c.logger.Log("level", "info", "message", "collected metrics", "scrapeID", c.scrapeID, "scrapeTime", scrapingElapsed.Seconds())
}

func (c *Collector) getNeighbours(n int, subsets []v1.EndpointSubset) ([]string, error) {
	// Find our IP - note: this does not open a connection, due to UDP.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer conn.Close()
	ip := conn.LocalAddr().(*net.UDPAddr).IP.String()

	// Get all other net-exporter IPs, and sort them.
	addresses := []string{}
	for _, endpointSubset := range subsets {
		for _, address := range endpointSubset.Addresses {
			addresses = append(addresses, address.IP)
		}
	}

	// Calculate n neighbours, given our local IP and all other net-exporter IPs.
	neighbours := c.calculateNeighbours(n, ip, addresses)

	c.logger.Log("level", "info", "message", "calculated neighbours", "ip", ip, "neighbours", strings.Join(neighbours, ", "), "scrapeID", c.scrapeID)

	return neighbours, nil
}

func (c *Collector) calculateNeighbours(n int, ip string, addresses []string) []string {
	sort.Strings(addresses)

	neighbours := []string{}

	if n > len(addresses) {
		n = len(addresses)
	}

	for i := 0; i < len(addresses); i++ {
		if addresses[i] == ip {
			for j := 1; j < n+1; j++ {
				k := i + j
				k = k % len(addresses)
				neighbours = append(neighbours, addresses[k])
			}
		}
	}

	return neighbours
}
