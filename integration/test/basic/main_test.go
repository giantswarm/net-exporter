// +build k8srequired

package basic

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giantswarm/apprclient"
	e2esetup "github.com/giantswarm/e2esetup/chart"
	"github.com/giantswarm/e2esetup/chart/env"
	"github.com/giantswarm/e2esetup/k8s"
	"github.com/giantswarm/e2etests/managedservices"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	testName = "basic"

	appName   = "net-exporter"
	chartName = "net-exporter"
)

var (
	a          *apprclient.Client
	ms         *managedservices.ManagedServices
	helmClient *helmclient.Client
	l          micrologger.Logger
	k8sSetup   *k8s.Setup
)

func init() {
	var err error

	{
		c := micrologger.Config{}

		l, err = micrologger.New(c)
		if err != nil {
			panic(err.Error())
		}
	}

	var cpK8sClients *k8s.Clients
	{
		c := k8s.ClientsConfig{
			Logger: l,

			KubeConfigPath: env.KubeConfigPath(),
		}

		cpK8sClients, err = k8s.NewClients(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := k8s.SetupConfig{
			Clients: cpK8sClients,
			Logger:  l,
		}

		k8sSetup, err = k8s.NewSetup(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := apprclient.Config{
			Fs:     afero.NewOsFs(),
			Logger: l,

			Address:      "https://quay.io",
			Organization: "giantswarm",
		}
		a, err = apprclient.New(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := helmclient.Config{
			Logger:     l,
			K8sClient:  cpK8sClients.K8sClient(),
			RestConfig: cpK8sClients.RestConfig(),

			TillerNamespace: "giantswarm",
		}

		helmClient, err = helmclient.New(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := managedservices.Config{
			ApprClient: a,
			Clients:    cpK8sClients,
			HelmClient: helmClient,
			Logger:     l,

			ChartConfig: managedservices.ChartConfig{
				ChannelName:     fmt.Sprintf("%s-%s", env.CircleSHA(), testName),
				ChartName:       chartName,
				ChartValues:     fmt.Sprintf("{ \"image\": { \"tag\": \"%s\" }, \"namespace\": \"%s\" }", env.CircleSHA(), metav1.NamespaceSystem),
				Namespace:       metav1.NamespaceSystem,
				RunReleaseTests: false,
			},
			ChartResources: managedservices.ChartResources{
				DaemonSets: []managedservices.DaemonSet{
					{
						Name:      appName,
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							"app": appName,
						},
						MatchLabels: map[string]string{
							"app": appName,
						},
					},
				},
			},
		}

		ms, err = managedservices.New(c)
		if err != nil {
			panic(err.Error())
		}
	}
}

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	ctx := context.Background()

	{
		c := e2esetup.Config{
			HelmClient: helmClient,
			Setup:      k8sSetup,
		}

		v, err := e2esetup.Setup(ctx, m, c)
		if err != nil {
			l.LogCtx(ctx, "level", "error", "message", "e2e test failed", "stack", microerror.Stack(err))
		}

		os.Exit(v)
	}
}
