package cilium

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
)

const (
	namespace = "cilium"

	bucketStart  = 0.001
	bucketFactor = 2
	numBuckets   = 10
)

// Config provides the necessary configuration for creating a Collector.
type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
}

// Collector implements the Collector interface, exposing Cilium functionality metrics.
type Collector struct {
	k8sClient  kubernetes.Interface
	logger     micrologger.Logger
	policyMaps *prometheus.GaugeVec
}

// New creates a Collector, given a Config.
func New(config Config) (*Collector, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	policyMaps := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: prometheus.BuildFQName(namespace, "", "policy_maps"),
		Help: "Number of policy maps.",
	}, []string{"endpoint_id"})

	prometheus.MustRegister(policyMaps)

	collector := &Collector{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		policyMaps: policyMaps,
	}

	return collector, nil
}

// Describe implements the Describe method of the Collector interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	// TODO: Add metric descriptions when implementing metrics
}

// Collect implements the Collect method of the Collector interface.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	// TODO: Implement metric collection
	policyMaps, _ := listAllMaps()
	for _, policyMap := range policyMaps {
		c.policyMaps.WithLabelValues(policyMap.EndpointID).Set(float64(policyMap.Size))
	}
}
