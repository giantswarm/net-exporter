package nic

import (
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/safchain/ethtool"
)

const (
	nic_metric_namespace = "nic"
)

// Config provides the necessary configuration for creating a Collector.
type Config struct {
	Logger micrologger.Logger

	IFace string
}

type Collector struct {
	logger micrologger.Logger

	iface   string
	metrics map[string]*prometheus.Desc
}

// New creates a Collector, given a Config.
func New(config Config) (*Collector, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.IFace == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.IFace must not be empty", config)
	}

	collector := &Collector{
		iface: config.IFace,
	}

	nicStats, err := ethtool.Stats(collector.iface)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	collector.metrics = make(map[string]*prometheus.Desc)
	for label, _ := range nicStats {
		fqName := prometheus.BuildFQName(nic_metric_namespace, "", label)
		collector.metrics[label] = prometheus.NewDesc(fqName, fmt.Sprintf("Generated description for metric %#q", label), []string{"hostname", "iface"}, nil)
	}

	return collector, nil
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {

	for _, nicMetric := range collector.metrics {
		ch <- nicMetric
	}
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	nicStats, _ := ethtool.Stats(collector.iface)

	for label, nicMetric := range collector.metrics {
		ch <- prometheus.MustNewConstMetric(nicMetric, prometheus.GaugeValue, float64(nicStats[label]), collector.iface)
	}
}
