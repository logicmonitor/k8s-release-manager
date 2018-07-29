package config

// Config represents the application's configuration
type Config struct {
	DebugMode   bool
	VerboseMode bool
	Helm        *HelmConfig
	Manager     *ManagerConfig
	Client      *ClientConfig
}

// HelmConfig represents the application's configurations for interacting with Helm
type HelmConfig struct {
	TillerHost      string
	TillerNamespace string `default:"kube-system"`
}

// ManagerConfig represents configurations for manager mode
type ManagerConfig struct {
	ReleaseName     string `default:""`
	PollingInterval int64
	StoragePath     string
}

// ClientConfig represents configurations for client mode
type ClientConfig struct {
	ReleaseTimeoutSec int64 `default:"300"`
}
