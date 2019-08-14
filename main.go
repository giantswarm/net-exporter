package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/giantswarm/exporterkit"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/net-exporter/dns"
	"github.com/giantswarm/net-exporter/network"
	"github.com/giantswarm/net-exporter/nic"
	"github.com/giantswarm/net-exporter/nstat"
)

var (
	nicExporter   bool
	nstatExporter bool
	hosts         string
	iface         string
	namespace     string
	port          string
	service       string
)

func init() {
	flag.BoolVar(&nicExporter, "nic-exporter-enabled", false, "nic exporter state (default value 'false')")
	flag.BoolVar(&nstatExporter, "nstat-exporter-enabled", false, "nstat exporter state (default value 'false')")
	flag.StringVar(&hosts, "hosts", "giantswarm.io,kubernetes.default.svc.cluster.local", "DNS hosts to resolve")
	flag.StringVar(&iface, "iface", "eth0", "Interface name to retrieve stats from")
	flag.StringVar(&namespace, "namespace", "monitoring", "Namespace of net-exporter service")
	flag.StringVar(&port, "port", "8000", "Port of net-exporter service")
	flag.StringVar(&service, "service", "net-exporter", "Name of net-exporter service")
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

	var collectors []prometheus.Collector

	var dnsCollector prometheus.Collector
	{
		splitHosts := strings.Split(hosts, ",")

		c := dns.Config{
			Logger: logger,

			Hosts: splitHosts,
		}

		dnsCollector, err = dns.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}

		collectors = append(collectors, dnsCollector)
	}

	var networkCollector prometheus.Collector
	{
		c := network.Config{
			Dialer: &net.Dialer{
				Timeout: 5 * time.Second,
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

		collectors = append(collectors, networkCollector)
	}

	if nicExporter {
		logger.Log("debug", "nic exporter enabled")

		c := nic.Config{
			Logger: logger,

			IFace: iface,
		}

		nicCollector, err := nic.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}

		collectors = append(collectors, nicCollector)
	}

	if nstatExporter {
		logger.Log("debug", "nstat exporter enabled")

		c := nstat.Config{
			Logger: logger,
		}

		nstatCollector, err := nstat.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}

		collectors = append(collectors, nstatCollector)
	}

	var exporter *exporterkit.Exporter
	{
		c := exporterkit.Config{
			Collectors: collectors,
			Logger:     logger,
		}

		exporter, err = exporterkit.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	exporter.Run()
}
