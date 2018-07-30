package config

// Config represents the application's configuration
type Config struct {
	Backend     *BackendConfig
	Client      *ClientConfig
	Export      *ExportConfig
	Helm        *HelmConfig
	DebugMode   bool
	DryRun      bool
	VerboseMode bool
}

// HelmConfig represents the application's configurations for interacting with Helm
type HelmConfig struct {
	TillerHost      string
	TillerNamespace string `default:"kube-system"`
}

//BackendConfig represents configuration options for the backend storage
type BackendConfig struct {
	StoragePath string
}

// ExportConfig represents configurations for manager mode
type ExportConfig struct {
	DaemonMode      bool
	ReleaseName     string
	PollingInterval int64
}

// ClientConfig represents configurations for client mode
type ClientConfig struct {
	ReleaseTimeoutSec int64 `default:"300"`
}
