package lmhelm

import (
	"fmt"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/portforwarder"
	rls "k8s.io/helm/pkg/proto/hapi/release"
)

// Client represents the LM helm client wrapper
type Client struct {
	Helm             *helm.Client
	helmConfig       *config.HelmConfig
	kubernetesClient *kubernetes.Clientset
	kubernetesConfig *rest.Config
	settings         helm_env.EnvSettings
}

// Init initializes the LM helm wrapper struct
func (c *Client) Init(helmConfig *config.HelmConfig, kubernetesClient *kubernetes.Clientset, kubernetesConfig *rest.Config) error {
	var err error
	c.helmConfig = helmConfig
	c.kubernetesClient = kubernetesClient
	c.kubernetesConfig = kubernetesConfig
	c.settings = c.getHelmSettings()

	c.Helm, err = c.newHelmClient()
	return err
}

// NewHeClient returns a helm client
func (c *Client) newHelmClient() (*helm.Client, error) {
	tillerHost, err := c.tillerHost()
	if err != nil {
		return nil, err
	}

	log.Debugf("Using tiller host %s", tillerHost)
	helmClient := helm.NewClient(helm.Host(tillerHost))
	return helmClient, nil
}

func (c *Client) tillerHost() (string, error) {
	log.Debugf("Setting up port forwarding tunnel to tiller")
	tunnel, err := portforwarder.New(c.settings.TillerNamespace, c.kubernetesClient, c.kubernetesConfig)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("127.0.0.1:%d", tunnel.Local), nil
}

// HelmSettings returns the helm client settings
func (c *Client) HelmSettings() helm_env.EnvSettings {
	return c.settings
}

// Config returns the client application settings
func (c *Client) Config() *config.HelmConfig {
	return c.helmConfig
}

func (c *Client) getHelmSettings() helm_env.EnvSettings {
	var settings helm_env.EnvSettings

	settings.TillerNamespace = c.helmConfig.TillerNamespace
	if settings.TillerNamespace == "" {
		settings.TillerNamespace = constants.DefaultTillerNamespace
	}
	return settings
}

// ListInstalledReleases lists all currently installed helm releases
func (c *Client) ListInstalledReleases() ([]*rls.Release, error) {
	rsp, err := c.Helm.ListReleases(listOpts()...)
	if err != nil {
		return nil, err
	}
	return rsp.Releases, nil
}

// Install a release
func (c *Client) Install(r *rls.Release) error {
	vals := []byte(r.GetConfig().GetRaw())
	log.Debugf("Installing release %s", r.GetName())
	rsp, err := c.Helm.InstallReleaseFromChart(r.GetChart(), r.GetNamespace(), installOpts(r, vals, c.helmConfig)...)
	if rsp != nil {
		log.Infof("Release %s status %s", rsp.Release.GetName(), rsp.Release.GetInfo().GetStatus().GetCode().String())
	} else {
		log.Errorf("%v", err)
	}
	return err
}
