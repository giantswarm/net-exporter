package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/giantswarm/exporterkit"
	"github.com/giantswarm/k8sclient/k8srestconfig"
	"github.com/giantswarm/micrologger"
	dnsclient "github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/net-exporter/dns"
	"github.com/giantswarm/net-exporter/network"
)

var (
	hosts     string
	namespace string
	port      string
	service   string
	timeout   time.Duration
)

func init() {
	flag.StringVar(&hosts, "hosts", "giantswarm.io.,kubernetes.default.svc.cluster.local.", "DNS hosts to resolve")
	flag.StringVar(&namespace, "namespace", "monitoring", "Namespace of net-exporter service")
	flag.StringVar(&port, "port", "8000", "Port of net-exporter service")
	flag.StringVar(&service, "service", "net-exporter", "Name of net-exporter service")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "Timeout of the dialer")
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--help") {
		return
	}

	flag.Parse()

	var err error

	var logger micrologger.Logger
	{
		logger, err = micrologger.New(micrologger.Config{})
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var restConfig *rest.Config
	{
		c := k8srestconfig.Config{
			Logger: logger,

			InCluster: true,
		}

		restConfig, err = k8srestconfig.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var kubernetesClient kubernetes.Interface
	{
		kubernetesClient, err = kubernetes.NewForConfig(restConfig)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var dnsCollector prometheus.Collector
	{
		splitHosts := strings.Split(hosts, ",")

		c := dns.Config{
			KubernetesClient: kubernetesClient,
			Logger:           logger,
			TCPClient: &dnsclient.Client{
				Net: "tcp",
			},
			UDPClient: &dnsclient.Client{
				Net: "udp",
			},

			Hosts: splitHosts,
		}

		dnsCollector, err = dns.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var networkCollector prometheus.Collector
	{
		c := network.Config{
			Dialer: &net.Dialer{
				Timeout: timeout,
			},
			KubernetesClient: kubernetesClient,
			Logger:           logger,

			Namespace: namespace,
			Port:      port,
			Service:   service,
		}

		networkCollector, err = network.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var exporter *exporterkit.Exporter
	{
		c := exporterkit.Config{
			Collectors: []prometheus.Collector{
				dnsCollector,
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
