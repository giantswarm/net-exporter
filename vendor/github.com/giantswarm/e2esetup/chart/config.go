package chart

import (
	"github.com/giantswarm/helmclient"
)

type Config struct {
	HelmClient *helmclient.Client
	Host       LegacyFramework
}
