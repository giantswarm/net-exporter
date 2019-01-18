package managedservices

import (
	"context"
	"fmt"
	"reflect"

	"github.com/giantswarm/apprclient"
	"github.com/giantswarm/e2e-harness/pkg/framework"
	frameworkresource "github.com/giantswarm/e2e-harness/pkg/framework/resource"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	ApprClient    *apprclient.Client
	HelmClient    *helmclient.Client
	HostFramework *framework.Host
	Logger        micrologger.Logger

	ChartConfig    ChartConfig
	ChartResources ChartResources
}

type ManagedServices struct {
	apprClient    *apprclient.Client
	helmClient    *helmclient.Client
	hostFramework *framework.Host
	logger        micrologger.Logger
	resource      *frameworkresource.Resource

	chartConfig    ChartConfig
	chartResources ChartResources
}

func New(config Config) (*ManagedServices, error) {
	var err error

	if config.ApprClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.ApprClient must not be empty", config)
	}
	if config.HostFramework == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HostFramework must not be empty", config)
	}
	if config.HelmClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HelmClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	err = config.ChartConfig.Validate()
	if err != nil {
		return nil, microerror.Mask(err)
	}
	err = config.ChartResources.Validate()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var resource *frameworkresource.Resource
	{
		c := frameworkresource.Config{
			ApprClient: config.ApprClient,
			HelmClient: config.HelmClient,
			Logger:     config.Logger,
			Namespace:  config.ChartConfig.Namespace,
		}

		resource, err = frameworkresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	ms := &ManagedServices{
		apprClient:    config.ApprClient,
		helmClient:    config.HelmClient,
		hostFramework: config.HostFramework,
		logger:        config.Logger,
		resource:      resource,

		chartConfig:    config.ChartConfig,
		chartResources: config.ChartResources,
	}

	return ms, nil
}

func (ms *ManagedServices) Test(ctx context.Context) error {
	var err error

	{
		ms.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("installing chart %#q", ms.chartConfig.ChartName))

		err = ms.resource.Install(ms.chartConfig.ChartName, ms.chartConfig.ChartValues, ms.chartConfig.ChannelName)
		if err != nil {
			return microerror.Mask(err)
		}

		ms.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("installed chart %#q", ms.chartConfig.ChartName))
	}

	{
		ms.logger.LogCtx(ctx, "level", "debug", "message", "waiting for deployed status")

		err = ms.resource.WaitForStatus(ms.chartConfig.ChartName, "DEPLOYED")
		if err != nil {
			return microerror.Mask(err)
		}

		ms.logger.LogCtx(ctx, "level", "debug", "message", "chart is deployed")
	}
	{
		ms.logger.LogCtx(ctx, "level", "debug", "message", "checking resources")

		for _, ds := range ms.chartResources.DaemonSets {
			ms.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("checking daemonset %#q", ds.Name))

			err = ms.checkDaemonSet(ds)
			if err != nil {
				return microerror.Mask(err)
			}

			ms.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("daemonset %#q is correct", ds.Name))
		}

		for _, d := range ms.chartResources.Deployments {
			ms.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("checking deployment %#q", d.Name))

			err = ms.checkDeployment(d)
			if err != nil {
				return microerror.Mask(err)
			}

			ms.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deployment %#q is correct", d.Name))
		}

		ms.logger.LogCtx(ctx, "level", "debug", "message", "resources are correct")
	}

	{
		if ms.chartConfig.RunReleaseTests {
			ms.logger.LogCtx(ctx, "level", "debug", "message", "running release tests")

			err = ms.helmClient.RunReleaseTest(ctx, ms.chartConfig.ChartName)
			if err != nil {
				return microerror.Mask(err)
			}

			ms.logger.LogCtx(ctx, "level", "debug", "message", "release tests passed")
		} else {
			ms.logger.LogCtx(ctx, "level", "debug", "message", "skipping release tests")
		}
	}

	return nil
}

// checkDaemonSet ensures that key properties of the daemonset are correct.
func (ms *ManagedServices) checkDaemonSet(expectedDaemonSet DaemonSet) error {
	ds, err := ms.hostFramework.K8sClient().Apps().DaemonSets(expectedDaemonSet.Namespace).Get(expectedDaemonSet.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return microerror.Maskf(notFoundError, "daemonset %#q", expectedDaemonSet.Name)
	} else if err != nil {
		return microerror.Mask(err)
	}

	err = ms.checkLabels("daemonset labels", expectedDaemonSet.Labels, ds.ObjectMeta.Labels)
	if err != nil {
		return microerror.Mask(err)
	}

	err = ms.checkLabels("daemonset matchLabels", expectedDaemonSet.MatchLabels, ds.Spec.Selector.MatchLabels)
	if err != nil {
		return microerror.Mask(err)
	}

	err = ms.checkLabels("daemonset pod labels", expectedDaemonSet.Labels, ds.Spec.Template.ObjectMeta.Labels)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// checkDeployment ensures that key properties of the deployment are correct.
func (ms *ManagedServices) checkDeployment(expectedDeployment Deployment) error {
	ds, err := ms.hostFramework.K8sClient().Apps().Deployments(expectedDeployment.Namespace).Get(expectedDeployment.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return microerror.Maskf(notFoundError, "deployment: %#q", expectedDeployment.Name)
	} else if err != nil {
		return microerror.Mask(err)
	}

	if int32(expectedDeployment.Replicas) != *ds.Spec.Replicas {
		return microerror.Maskf(invalidReplicasError, "expected %d replicas got: %d", expectedDeployment.Replicas, *ds.Spec.Replicas)
	}

	err = ms.checkLabels("deployment labels", expectedDeployment.DeploymentLabels, ds.ObjectMeta.Labels)
	if err != nil {
		return microerror.Mask(err)
	}

	err = ms.checkLabels("deployment matchLabels", expectedDeployment.MatchLabels, ds.Spec.Selector.MatchLabels)
	if err != nil {
		return microerror.Mask(err)
	}

	err = ms.checkLabels("deployment pod labels", expectedDeployment.PodLabels, ds.Spec.Template.ObjectMeta.Labels)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (ms *ManagedServices) checkLabels(labelType string, expectedLabels, labels map[string]string) error {
	if !reflect.DeepEqual(expectedLabels, labels) {
		ms.logger.Log("level", "debug", "message", fmt.Sprintf("expected %s: %v got: %v", labelType, expectedLabels, labels))
		return microerror.Maskf(invalidLabelsError, "%s do not match expected labels", labelType)
	}

	return nil
}
