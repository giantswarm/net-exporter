package nstat

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	nstat_metric_namespace = "nstat"
)

// Config provides the necessary configuration for creating a Collector.
type Config struct {
	Logger micrologger.Logger
}

type Collector struct {
	metrics map[string]*prometheus.Desc
}

// New creates a Collector, given a Config.
func New(config Config) (*Collector, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	nstatOutput, err := exec.Command("/sbin/nstat", "-a", "--json").Output()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var nstatObj map[string]interface{}

	err = json.Unmarshal([]byte(nstatOutput), &nstatObj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	collector := &Collector{}
	metrics, _ := nstatObj["kernel"].(map[string]interface{})

	collector.metrics = make(map[string]*prometheus.Desc)
	for label, _ := range metrics {
		fqName := prometheus.BuildFQName(nstat_metric_namespace, "", label)
		collector.metrics[label] = prometheus.NewDesc(fqName, fmt.Sprintf("Generated description for metric %#q", label), []string{}, nil)
	}

	return collector, nil
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {

	for _, nstatMetric := range collector.metrics {
		ch <- nstatMetric
	}
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	nstatOutput, err := exec.Command("/sbin/nstat", "-a", "--json").Output()
	if err != nil {
		log.Fatal(err)
	}

	var nstatObj map[string]interface{}

	err = json.Unmarshal([]byte(nstatOutput), &nstatObj)
	if err != nil {
		log.Fatal(err)
	}

	kernelMetrics := nstatObj["kernel"].(map[string]interface{})

	for label, nstatMetric := range collector.metrics {
		ch <- prometheus.MustNewConstMetric(nstatMetric, prometheus.GaugeValue, kernelMetrics[label].(float64))
	}
}
