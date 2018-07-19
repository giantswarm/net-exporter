package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/giantswarm/exporterkit"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/net-exporter/network"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--help") {
		return
	}

	var err error

	var logger micrologger.Logger
	{
		logger, err = micrologger.New(micrologger.Config{})
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var kubernetesClient kubernetes.Interface
	{
		var config *rest.Config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}

		kubernetesClient, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var networkCollector prometheus.Collector
	{
		c := network.Config{
			Dialer: &net.Dialer{
				Timeout: 5 * time.Second,
			},
			KubernetesClient: kubernetesClient,
			Logger:           logger,

			Host:        "net-exporter:8000",
			Namespace:   "monitoring",
			Port:        "8000",
			ServiceName: "net-exporter",
		}

		networkCollector, err = network.NewCollector(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var exporter *exporterkit.Exporter
	{
		c := exporterkit.Config{
			Collectors: []prometheus.Collector{
				networkCollector,
			},
			Logger: logger,
		}

		exporter, err = exporterkit.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	exporter.Run()
}
