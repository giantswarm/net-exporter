package chart

import (
	"github.com/giantswarm/helmclient"

	"github.com/giantswarm/e2esetup/k8s"
)

type Config struct {
	HelmClient *helmclient.Client
	Setup      *k8s.Setup
}
