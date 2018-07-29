package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents the application's configuration file.
type Config struct {
	DebugMode          bool
	ManagerReleaseName string `default:""`
	PollingInterval    int64
	ReleaseTimeoutSec  int64 `default:"300"`
	StoragePath        string
	TillerHost         string
	TillerNamespace    string `default:"kube-system"`
	VerboseMode        bool
}

// New returns the application configuration specified by the config file.
func New() (*Config, error) {
	c := &Config{}
	if err := envconfig.Process("rlsmgr", c); err != nil {
		return nil, err
	}

	return c, nil
}
