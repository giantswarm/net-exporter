package chart

import (
	"context"
	"testing"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/e2esetup/chart/env"
)

func Setup(ctx context.Context, m *testing.M, config Config) (int, error) {
	var v int
	var err error
	var errors []error

	if config.HelmClient == nil {
		return v, microerror.Maskf(invalidConfigError, "%T.HelmClient must not be empty", config)
	}
	if config.Setup == nil {
		return v, microerror.Maskf(invalidConfigError, "%T.Setup must not be empty", config)
	}

	err = config.Setup.EnsureNamespaceCreated(ctx, "giantswarm")
	if err != nil {
		errors = append(errors, err)
		v = 1
	}

	err = config.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		errors = append(errors, err)
		v = 1
	}

	if v == 0 {
		v = m.Run()
	}

	if env.KeepResources() != "true" {
		// Only do full teardown when not on CI.
		if env.CircleCI() != "true" {
			err := teardown(config)
			if err != nil {
				errors = append(errors, err)
				v = 1
			}
		}
	}

	if len(errors) > 0 {
		return v, microerror.Mask(errors[0])
	}

	return v, nil
}
