package lmhelm

import (
	"fmt"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
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
	Helm         *helm.Client
	rlsmgrconfig *config.Config
	restConfig   *rest.Config
	settings     helm_env.EnvSettings
}

// Init initializes the LM helm wrapper struct
func (c *Client) Init(rlsmgrconfig *config.Config, config *rest.Config) error {
	// Instantiate the Helm client
	c.rlsmgrconfig = rlsmgrconfig
	c.settings = c.getHelmSettings()
	c.restConfig = config

	var err error
	c.Helm, err = c.newHelmClient()
	return err
}

// NewHeClient returns a helm client
func (c *Client) newHelmClient() (*helm.Client, error) {
	tillerHost, err := c.tillerHost()
	if err != nil {
		return nil, err
	}

	log.Infof("Using tiller host %s", tillerHost)
	heClient := helm.NewClient(helm.Host(tillerHost))
	return heClient, nil
}

func (c *Client) tillerHost() (string, error) {
	if c.settings.TillerHost != "" {
		return c.settings.TillerHost, nil
	}

	log.Debugf("Creating kubernetes client")
	client, err := kubernetes.NewForConfig(c.restConfig)
	if err != nil {
		return "", err
	}
	log.Debugf("Created kubernetes client")

	log.Debugf("Setting up port forwarding tunnel to tiller")
	tunnel, err := portforwarder.New(c.settings.TillerNamespace, client, c.restConfig)
	if err != nil {
		return "", err
	}
	log.Debugf("Set up port forwarding tunnel to tiller")

	return fmt.Sprintf("127.0.0.1:%d", tunnel.Local), nil
}

// HelmSettings returns the helm client settings
func (c *Client) HelmSettings() helm_env.EnvSettings {
	return c.settings
}

// Config returns the client application settings
func (c *Client) Config() *config.Config {
	return c.rlsmgrconfig
}

func (c *Client) getHelmSettings() helm_env.EnvSettings {
	var settings helm_env.EnvSettings
	settings.TillerHost = c.rlsmgrconfig.TillerHost
	settings.TillerNamespace = c.rlsmgrconfig.TillerNamespace
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

// func helmInstall(r *Release vals []byte) (*rls.Release, error) {
// 	log.Infof("Installing release %s", r.Name())
// 	rsp, err := r.Client.Helm.InstallReleaseFromChart(chart, r.Chartmgr.ObjectMeta.Namespace, installOpts(r, vals)...)
// 	if rsp == nil || rsp.Release == nil {
// 		rls, _ := getInstalledRelease(r)
// 		if rls != nil {
// 			return rls, nil
// 		}
// 	} else {
// 		return rsp.Release, nil
// 	}
// 	return nil, err
// }

// func helmUpdate(r *Release, chart *chart.Chart, vals []byte) (*rls.Release, error) {
// 	log.Infof("Updating release %s", r.Name())
// 	rsp, err := r.Client.Helm.UpdateReleaseFromChart(r.Name(), chart, updateOpts(r, vals)...)
// 	if rsp == nil || rsp.Release == nil {
// 		rls, _ := getInstalledRelease(r)
// 		if rls != nil {
// 			return rls, nil
// 		}
// 	} else {
// 		return rsp.Release, nil
// 	}
// 	return nil, err
// }
