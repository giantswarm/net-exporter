package chart

import (
	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/helmclient"
)

type Config struct {
	HelmClient *helmclient.Client
	Host       *framework.Host
}
