package config

// Config represents the application's configuration
// codebeat:disable[TOO_MANY_IVARS]
type Config struct {
	Backend       *BackendConfig
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
	ReleaseTimeoutSec int64
	TillerNamespace   string
}

//TransferConfig represents configuration options for the backend storage
type TransferConfig struct {
	Force          bool
	NewStoragePath string
}
