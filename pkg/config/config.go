package config

import "time"

// Config represents the application's configuration
// codebeat:disable[TOO_MANY_IVARS]
type Config struct {
	Backend       *BackendConfig
	Export        *ExportConfig
	ClusterConfig *ClusterConfig
	Import        *ImportConfig
	OptionsConfig OptionsConfig
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
	Namespaces      []string
}

//ImportConfig represents configuration options for the backend storage
type ImportConfig struct {
	Force             bool
	NewStoragePath    string
	Namespace         string
	Target            string
	Values            map[string]string
	ExcludeNamespaces []string
	Threads           int64
}

// OptionsConfig represents the client configurations options for listing and installing releases
type OptionsConfig struct {
	Install *InstallConfig
	List    *ListConfig
}

// InstallConfig represents the configuration options for installing releases
type InstallConfig struct {
	Wait            bool
	Timeout         time.Duration
	Replace         bool
	CreateNamespace bool
	Atomic          bool
	DryRun          bool
}

// ListConfig represents the configuration options for listing releases
type ListConfig struct {
	Deployed      bool
	Failed        bool
	AllNamespaces bool
}
