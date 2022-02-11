package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/giantswarm/exporterkit"
	"github.com/giantswarm/k8sclient/v4/pkg/k8srestconfig"
	"github.com/giantswarm/micrologger"
	dnsclient "github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/net-exporter/dns"
	"github.com/giantswarm/net-exporter/network"
	"github.com/giantswarm/net-exporter/ntp"
)

var (
	disableDNSTCPCheck bool
	hosts              string
	dnsService         string
	namespace          string
	ntpServers         string
	port               string
	service            string
	timeout            time.Duration
)

func init() {
	flag.BoolVar(&disableDNSTCPCheck, "disable-dns-tcp-check", false, "Disable DNS TCP check")
	flag.StringVar(&hosts, "hosts", "giantswarm.io.,kubernetes.default.svc.cluster.local.", "DNS hosts to resolve")
	flag.StringVar(&dnsService, "dns-service", "coredns", "Name of DNS service")
	flag.StringVar(&namespace, "namespace", "monitoring", "Namespace of net-exporter service")
	flag.StringVar(&ntpServers, "ntp-servers", "0.flatcar.pool.ntp.org,1.flatcar.pool.ntp.org", "NTP servers to use for time synchronization")
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

	var k8sClient kubernetes.Interface
	{
		k8sClient, err = kubernetes.NewForConfig(restConfig)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var dnsCollector prometheus.Collector
	{
		splitHosts := strings.Split(hosts, ",")

		c := dns.Config{
			K8sClient: k8sClient,
			Logger:    logger,
			TCPClient: &dnsclient.Client{
				Net: "tcp",
			},
			UDPClient: &dnsclient.Client{
				Net: "udp",
			},

			DisableTCPCheck: disableDNSTCPCheck,
			Hosts:           splitHosts,
			Service:         dnsService,
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
			K8sClient: k8sClient,
			Logger:    logger,

			Namespace: namespace,
			Port:      port,
			Service:   service,
		}

		networkCollector, err = network.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var ntpCollector prometheus.Collector
	{
		splitNTPServers := strings.Split(ntpServers, ",")

		c := ntp.Config{
			Logger: logger,

			NTPServers: splitNTPServers,
		}

		ntpCollector, err = ntp.New(c)
		if err != nil {
			panic(microerror.JSON(err))
		}
	}

	var exporter *exporterkit.Exporter
	{
		c := exporterkit.Config{
			Address: "http://0.0.0.0:" + port,
			Collectors: []prometheus.Collector{
				dnsCollector,
				networkCollector,
				ntpCollector,
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
