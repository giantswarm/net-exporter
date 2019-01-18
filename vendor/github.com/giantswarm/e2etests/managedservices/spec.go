package managedservices

import (
	"context"

	"github.com/giantswarm/microerror"
)

// ChartConfig is the chart to test.
type ChartConfig struct {
	ChannelName     string
	ChartName       string
	ChartValues     string
	Namespace       string
	RunReleaseTests bool
}

func (cc ChartConfig) Validate() error {
	if cc.ChannelName == "" {
		return microerror.Maskf(invalidConfigError, "%T.ChannelName must not be empty", cc)
	}
	if cc.ChartName == "" {
		return microerror.Maskf(invalidConfigError, "%T.ChartName must not be empty", cc)
	}
	if cc.Namespace == "" {
		return microerror.Maskf(invalidConfigError, "%T.Namespace must not be empty", cc)
	}

	return nil
}

// ChartResources are the key resources deployed by the chart.
type ChartResources struct {
	DaemonSets  []DaemonSet
	Deployments []Deployment
}

func (cr ChartResources) Validate() error {
	if len(cr.DaemonSets) == 0 && len(cr.Deployments) == 0 {
		return microerror.Maskf(invalidConfigError, "at least one daemonset or deployment must be specified")
	}

	return nil
}

// DaemonSet is a daemonset to be tested.
type DaemonSet struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	MatchLabels map[string]string
	Replicas    int
}

// Deployment is a deployment to be tested.
type Deployment struct {
	Name             string
	Namespace        string
	DeploymentLabels map[string]string
	MatchLabels      map[string]string
	PodLabels        map[string]string
	Replicas         int
}

type Interface interface {
	// Test executes the test of a managed services chart with basic
	// functionality that applies to all managed services charts.
	//
	// - Install chart.
	// - Check chart is deployed.
	// - Check key resources are correct.
	// - Run helm release tests if configured.
	//
	Test(ctx context.Context) error
}
