package network

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/giantswarm/exporterkit"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	Namespace = "network"
)

type Config struct {
	Dialer           *net.Dialer
	KubernetesClient kubernetes.Interface
	Logger           micrologger.Logger

	Host        string
	Namespace   string
	Port        string
	ServiceName string
}

type NetworkCollector struct {
	dialer           *net.Dialer
	kubernetesClient kubernetes.Interface
	logger           micrologger.Logger

	host        string
	namespace   string
	port        string
	serviceName string

	count                map[string]float64
	errorCount           map[string]float64
	kubernetesErrorCount float64

	countMutex      sync.Mutex
	errorCountMutex sync.Mutex

	total                *prometheus.Desc
	errorTotal           *prometheus.Desc
	latency              *prometheus.Desc
	kubernetesErrorTotal *prometheus.Desc
}

func NewCollector(config Config) (*NetworkCollector, error) {
	if config.Dialer == nil {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.Dialer must not be empty", config)
	}
	if config.KubernetesClient == nil {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.KubernetesClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.Host == "" {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.Host must not be empty", config)
	}
	if config.Namespace == "" {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.Namespace must not be empty", config)
	}
	if config.Port == "" {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.Port must not be empty", config)
	}
	if config.ServiceName == "" {
		return nil, microerror.Maskf(exporterkit.InvalidConfigError, "%T.ServiceName must not be empty", config)
	}

	networkCollector := NetworkCollector{
		dialer:           config.Dialer,
		kubernetesClient: config.KubernetesClient,
		logger:           config.Logger,

		host:        config.Host,
		namespace:   config.Namespace,
		port:        config.Port,
		serviceName: config.ServiceName,

		count:      map[string]float64{},
		errorCount: map[string]float64{},

		total: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "check_total"),
			"Total number of network checks.",
			[]string{"host"},
			nil,
		),
		errorTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "check_error_total"),
			"Total number of network check errors.",
			[]string{"host"},
			nil,
		),
		latency: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "check_seconds"),
			"Time taken to successfully check network.",
			[]string{"host"},
			nil,
		),
		kubernetesErrorTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "kubernetes_error_total"),
			"Total number of errors reaching Kubernetes API.",
			nil,
			nil,
		),
	}

	return &networkCollector, nil
}

func (c *NetworkCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.errorTotal
	ch <- c.latency
	ch <- c.kubernetesErrorTotal
}

func (c *NetworkCollector) Collect(ch chan<- prometheus.Metric) {
	hosts := []string{c.host}

	endpoints, err := c.kubernetesClient.CoreV1().Endpoints(c.namespace).Get(c.serviceName, metav1.GetOptions{})
	if err != nil {
		c.logger.Log("level", "error", "message", "could not get endpoints", "stack", fmt.Sprintf("%#v", err))
		c.kubernetesErrorCount++
		ch <- prometheus.MustNewConstMetric(c.kubernetesErrorTotal, prometheus.CounterValue, c.kubernetesErrorCount)

		return
	}

	for _, endpointSubset := range endpoints.Subsets {
		for _, address := range endpointSubset.Addresses {
			hosts = append(hosts, fmt.Sprintf("%v:%v", address.IP, c.port))
		}
	}

	var wg sync.WaitGroup

	for _, host := range hosts {
		wg.Add(1)

		go func(host string) {
			defer wg.Done()

			start := time.Now()

			c.countMutex.Lock()
			c.count[host]++
			c.countMutex.Unlock()

			conn, err := c.dialer.Dial("tcp", host)
			if err != nil {
				c.logger.Log("level", "error", "message", "could not dial host", "host", host, "stack", fmt.Sprintf("%#v", err))

				c.errorCountMutex.Lock()
				c.errorCount[host]++
				ch <- prometheus.MustNewConstMetric(c.errorTotal, prometheus.CounterValue, c.errorCount[host], host)
				c.errorCountMutex.Unlock()

				return
			}
			defer conn.Close()

			elapsed := time.Since(start)

			c.countMutex.Lock()
			ch <- prometheus.MustNewConstMetric(c.total, prometheus.CounterValue, c.count[host], host)
			c.countMutex.Unlock()
			ch <- prometheus.MustNewConstMetric(c.latency, prometheus.GaugeValue, elapsed.Seconds(), host)
		}(host)
	}

	wg.Wait()
}
