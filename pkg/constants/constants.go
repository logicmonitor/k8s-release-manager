package constants

var (
	// Version is the Release Manager version and is set at build time.
	Version string
)

const (
	// HelmStableRepo is the name of the stable helm repo
	HelmStableRepo = "stable"
	// HelmStableRepoURL is the URL of the stable helm repo
	HelmStableRepoURL = "https://kubernetes-charts.storage.googleapis.com"
)

const (
	// RlsMgrSecretName is the service account name with the proper RBAC policies to allow a rlsmgr to poll Tiller.
	RlsMgrSecretName = "rlsmgr"
	//ManagerStateFilename is the filename used to store the manager state in the backend
	ManagerStateFilename = "rlsmgrstate.json"
	// ReleaseExtension is the file extension to use when storing releases in the backend
	ReleaseExtension = "release"
)

const (
	// DefaultKubeConfig file name
	DefaultKubeConfig = "config"
	// DefaultKubeConfigDir name
	DefaultKubeConfigDir = ".kube"
	// DefaultTillerNamespace kube-system
	DefaultTillerNamespace = "kube-system"
	//EnvKubeConfig is the default KUBECONFIG env var
	EnvKubeConfig = "KUBECONFIG"
)

const (
	// ValueStoragePath is the helm --set path for --path
	ValueStoragePath = "backend.storagePath"
)
