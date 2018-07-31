package config

// Config represents the application's configuration
type Config struct {
	Backend       *BackendConfig
	Client        *ClientConfig
	Export        *ExportConfig
	Helm          *HelmConfig
	ClusterConfig *ClusterConfig
	Transfer      *TransferConfig
	DebugMode     bool
	DryRun        bool
	VerboseMode   bool
}

//BackendConfig represents configuration options for the backend storage
type BackendConfig struct {
	StoragePath string
}

// ClientConfig represents configurations for client mode
type ClientConfig struct {
	ReleaseTimeoutSec int64 `default:"300"`
}

//ClusterConfig represents kubernetes configuration options
type ClusterConfig struct {
	KubeConfig  string
	KubeContext string
}

// ExportConfig represents configurations for manager mode
type ExportConfig struct {
	DaemonMode      bool
	ReleaseName     string
	PollingInterval int64
}

// HelmConfig represents the application's configurations for interacting with Helm
type HelmConfig struct {
	TillerNamespace string
}

}
