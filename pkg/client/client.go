package client

import (
	"strings"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	"github.com/logicmonitor/k8s-release-manager/pkg/utilities"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	// load kubernetes auth method
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesClient returns a kubernetes clientset from the given configuration
func KubernetesClient(c *config.ClusterConfig) (*kubernetes.Clientset, *rest.Config, error) {
	clusterConfig, err := clusterConfig(c)
	if err != nil {
		return nil, nil, err
	}

	log.Debugf("Creating kubernetes client")
	client, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, nil, err
	}
	return client, clusterConfig, nil
}

func clusterConfig(c *config.ClusterConfig) (*rest.Config, error) {
	kubeConfig := kubeConfig(c)
	// if kubeconfig set, use client config
	if kubeConfig != "" {
		c.KubeConfig = kubeConfig
		return localConfig(c)
	}

	// else, use in cluster config
	log.Debugf("Using in-cluster config")
	return rest.InClusterConfig()
}

func localConfig(c *config.ClusterConfig) (*rest.Config, error) {

	// rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules := &clientcmd.ClientConfigLoadingRules{
		DefaultClientConfig: &clientcmd.DefaultClientConfig,
	}
	overrides := &clientcmd.ConfigOverrides{}

	if c.KubeContext != "" {
		log.Debugf("Using kube context %s", c.KubeContext)
		overrides.CurrentContext = c.KubeContext
	}

	log.Debugf("Using kubeconfig %s", c.KubeConfig)
	rules.ExplicitPath = c.KubeConfig

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
}

func kubeConfig(c *config.ClusterConfig) string {
	if c.KubeConfig != "" {
		path, err := homedir.Expand(c.KubeConfig)
		if err != nil {
			log.Warnf("Unable to expand home directory")
			return c.KubeConfig
		}
		return path
	}

	defaultPath := defaultKubeConfigPath()
	if utilities.FileExists(defaultPath) {
		return defaultPath
	}
	return ""
}

func defaultKubeConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Errorf("Unable to determine home directory")
		return ""
	}
	return strings.Join([]string{home, constants.DefaultKubeConfigDir, constants.DefaultKubeConfig}, "/")
}
