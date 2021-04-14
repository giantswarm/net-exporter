module github.com/giantswarm/net-exporter

go 1.14

require (
	github.com/beevik/ntp v0.3.0
	github.com/giantswarm/exporterkit v0.2.1
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/micrologger v0.5.0
	github.com/google/go-cmp v0.5.5
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/miekg/dns v1.1.40
	github.com/prometheus/client_golang v1.9.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	k8s.io/api v0.18.5
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.5
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19 // indirect
)

replace (
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.25+incompatible
)
