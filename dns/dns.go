package dns

import (
	"fmt"
	"sync"
	"time"

	"github.com/giantswarm/exporterkit/histogramvec"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	dnsclient "github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	namespace = "dns"

	bucketStart  = 0.001
	bucketFactor = 2
	numBuckets   = 15
)

// Config provides the necessary configuration for creating a Collector.
type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
	TCPClient *dnsclient.Client
	UDPClient *dnsclient.Client

	DisableTCPCheck bool
	Hosts           []string
}

// Collector implements the Collector interface, exposing DNS latency information.
type Collector struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger
	tcpClient *dnsclient.Client
	udpClient *dnsclient.Client

	disableTCPCheck bool
	hosts           []string

	tcpLatencyHistogramVec  *histogramvec.HistogramVec
	tcpLatencyHistogramDesc *prometheus.Desc
	udpLatencyHistogramVec  *histogramvec.HistogramVec
	udpLatencyHistogramDesc *prometheus.Desc

	errorCount        prometheus.Counter
	resolveErrorCount *prometheus.CounterVec
}

// New creates a Collector, given a Config.
func New(config Config) (*Collector, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.TCPClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.TCPClient must not be empty", config)
	}
	if config.UDPClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.UDPClient must not be empty", config)
	}

	if len(config.Hosts) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.Host must not be empty", config)
	}

	var err error
	var tcpLatencyHistogramVec *histogramvec.HistogramVec
	{
		c := histogramvec.Config{
			BucketLimits: prometheus.ExponentialBuckets(bucketStart, bucketFactor, numBuckets),
		}
		tcpLatencyHistogramVec, err = histogramvec.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	var udpLatencyHistogramVec *histogramvec.HistogramVec
	{
		c := histogramvec.Config{
			BucketLimits: prometheus.ExponentialBuckets(bucketStart, bucketFactor, numBuckets),
		}
		udpLatencyHistogramVec, err = histogramvec.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	errorCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "error_total"),
		Help: "Total number of internal errors.",
	})
	resolveErrorCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, "", "resolve_error_total"),
			Help: "Total number of errors resolving hosts.",
		},
		[]string{"proto", "host"},
	)

	prometheus.MustRegister(errorCount)
	prometheus.MustRegister(resolveErrorCount)

	collector := &Collector{
		k8sClient: config.K8sClient,
		logger:    config.Logger,
		tcpClient: config.TCPClient,
		udpClient: config.UDPClient,

		disableTCPCheck: config.DisableTCPCheck,
		hosts:           config.Hosts,

		tcpLatencyHistogramVec: tcpLatencyHistogramVec,
		tcpLatencyHistogramDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "tcp_latency_seconds"),
			"Histogram of latency of TCP DNS resolutions.",
			[]string{"host"},
			nil,
		),
		udpLatencyHistogramVec: udpLatencyHistogramVec,
		udpLatencyHistogramDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "udp_latency_seconds"),
			"Histogram of latency of UDP DNS resolutions.",
			[]string{"host"},
			nil,
		),

		errorCount:        errorCount,
		resolveErrorCount: resolveErrorCount,
	}

	return collector, nil
}

// Describe implements the Describe method of the Collector interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	if !c.disableTCPCheck {
		ch <- c.tcpLatencyHistogramDesc
	}
	ch <- c.udpLatencyHistogramDesc
}

func (c *Collector) resolve(proto string, client *dnsclient.Client, host string, dnsServer string, latencyHistogramVec *histogramvec.HistogramVec) {
	start := time.Now()

	message := &dnsclient.Msg{}
	message.SetQuestion(host, dnsclient.TypeA)

	msg, _, err := client.Exchange(message, fmt.Sprintf("%s:53", dnsServer))
	if err != nil || len(msg.Answer) == 0 {
		c.logger.Log("level", "error", "message", fmt.Sprintf("could not resolve dns for host %#q and protocol %#q", host, proto), "stack", microerror.JSON(err))
		c.resolveErrorCount.WithLabelValues(proto, host).Inc()
		return
	}

	elapsed := time.Since(start)

	err = latencyHistogramVec.Add(host, elapsed.Seconds())
	if err != nil {
		c.logger.Log("level", "error", "message", fmt.Sprintf("failed to update latency histogram for host %#q and protocol %#q", host, proto), "stack", microerror.JSON(err))
		c.resolveErrorCount.WithLabelValues(proto, host).Inc()
		return
	}
}

// Collect implements the Collect method of the Collector interface.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	service, err := c.k8sClient.CoreV1().Services("kube-system").Get("coredns", metav1.GetOptions{})
	if err != nil {
		c.logger.Log("level", "error", "message", "could not collect service from kubernetes api", "stack", microerror.JSON(err))
		c.errorCount.Inc()
		return
	}

	var wg sync.WaitGroup

	for _, host := range c.hosts {
		if !c.disableTCPCheck {
			wg.Add(1)
			go func(host string) {
				defer wg.Done()

				c.resolve("tcp", c.tcpClient, host, service.Spec.ClusterIP, c.tcpLatencyHistogramVec)
			}(host)
		}

		wg.Add(1)
		go func(host string) {
			defer wg.Done()

			c.resolve("udp", c.udpClient, host, service.Spec.ClusterIP, c.udpLatencyHistogramVec)
		}(host)
	}

	wg.Wait()

	c.tcpLatencyHistogramVec.Ensure(c.hosts)
	c.udpLatencyHistogramVec.Ensure(c.hosts)

	if !c.disableTCPCheck {
		for host, histogram := range c.tcpLatencyHistogramVec.Histograms() {
			ch <- prometheus.MustNewConstHistogram(
				c.tcpLatencyHistogramDesc,
				histogram.Count(), histogram.Sum(), histogram.Buckets(),
				host,
			)
		}
	}
	for host, histogram := range c.udpLatencyHistogramVec.Histograms() {
		ch <- prometheus.MustNewConstHistogram(
			c.udpLatencyHistogramDesc,
			histogram.Count(), histogram.Sum(), histogram.Buckets(),
			host,
		)
	}
}
