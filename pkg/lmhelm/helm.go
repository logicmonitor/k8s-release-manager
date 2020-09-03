package lmhelm

import (
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/kube"
	rls "helm.sh/helm/v3/pkg/release"
)

// Client represents the LM helm client v3 wrapper
type Client struct {
	helmConfig    *action.Configuration
	settings      *cli.EnvSettings
	clusterConfig *config.ClusterConfig
	optionsConfig config.OptionsConfig
}

// Init initializes the LM helm wrapper struct
func (c *Client) Init(clusterConfig *config.ClusterConfig, optionsConfig config.OptionsConfig) error {

	var err error
	c.settings = c.getHelmSettings()
	c.clusterConfig = clusterConfig
	c.optionsConfig = optionsConfig

	return err
}

func (c *Client) getHelmSettings() *cli.EnvSettings {
	return cli.New()
}

// ListInstalledReleases lists all currently installed helm releases
func (c *Client) ListInstalledReleases() ([]*rls.Release, error) {

	if err := c.initActionConfig(""); err != nil {
		return nil, err
	}

	list := action.NewList(c.helmConfig)

	// List Options:
	list.Deployed = c.optionsConfig.List.Deployed
	list.Failed = c.optionsConfig.List.Failed
	list.AllNamespaces = c.optionsConfig.List.AllNamespaces

	results, err := list.Run()
	if err != nil {
		return nil, err
	}

	return results, nil

}

// Install ...
func (c *Client) Install(r *rls.Release) error {

	if err := c.initActionConfig(r.Namespace); err != nil {
		return err
	}

	install := action.NewInstall(c.helmConfig)

	// Install Options:
	install.Wait = c.optionsConfig.Install.Wait
	install.Timeout = c.optionsConfig.Install.Timeout
	install.Replace = c.optionsConfig.Install.Replace
	install.CreateNamespace = c.optionsConfig.Install.CreateNamespace
	install.Atomic = c.optionsConfig.Install.Atomic
	install.DryRun = c.optionsConfig.Install.DryRun

	install.ReleaseName = r.Name
	install.Namespace = r.Namespace

	log.Debugf("Installing release %s", r.Name)

	rsp, err := install.Run(r.Chart, r.Config)

	if rsp != nil {
		log.Infof("Release %s status %s", rsp.Name, rsp.Info.Status.String())
	}
	return err

}
func (c *Client) initActionConfig(namespace string) error {

	var err error

	if c.clusterConfig.KubeConfig == "" {
		c.helmConfig, err = getActionConfig(c.settings.KubeConfig, c.settings.KubeContext, namespace)
		if err != nil {
			return err
		}
	} else {
		c.helmConfig, err = getActionConfig(c.clusterConfig.KubeConfig, c.clusterConfig.KubeContext, namespace)
		if err != nil {
			return err
		}
	}
	return err
}

func getActionConfig(kubeConfig, context, namespace string) (*action.Configuration, error) {

	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(kube.GetConfig(kubeConfig, context, namespace), namespace, constants.HelmDriver, log.Printf); err != nil {
		return nil, err
	}

	return actionConfig, nil
}
