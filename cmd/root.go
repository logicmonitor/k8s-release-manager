package cmd

import (
	"fmt"
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rlsmgrconfig *config.Config
var debug bool
var dryRun bool
var verbose bool
var kubeConfig string
var kubeContext string
var storagePath string
var releaseName string
var tillerNamespace string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "releasemanager",
	Short: `
Release Manager is a tool for managing Helm release state`,
	Long: `
Release Manager provides functionality for exporting up and importing the
state of Helm releases currently deployed to a Kubernetes cluster. The state
of the installed releases is saved to a configurable backend for easy
restoration of previously-deployed releases or for simplified re-deployment of
those releases to a new Kubernetes cluster.

Release Manager operations can be run locally or within the Kubernetes cluster.
The application also supports a daemon mode that will periodically update the
saved state.

To export releases, Release Manager queries Tiller to collect metadata for all
releases currently deployed in the source cluster and writes this metadata to
the configured backend data store. If the Release Manager is deployed in
daemon mode via its own Helm chart, it will also store metadata about itself.
This metadata is used to prevent import operations from creating a new Release
Manager with the same configuration as the previous managed, causing both
instances to write conflicting state to the backend.

To import releases, Release Manager retrieves the state stored in the backend,
connects to Tiller in the target Kubernetes cluster, and deploys the saved
releases to the cluster.

Release Manager will use --kubeconfig/--kubecontext, $KUBECONFIG, or
~/.kube/config to establish a connection to the Kubernetes cluster. If none of
these configuraitons are set, an in-cluster connection will be attempted. All
actions will be performed against the current cluster and a given command will
only perform actions against a single cluster, i.e. 'export' will
export releases from the configured cluster while 'import' will deploy releases
to the configured cluster and 'clear' requires no custer connection whatsoever.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
		if dryRun {
			fmt.Println("Dry run. No changes will be made.")
		}
		rlsmgrconfig.DebugMode = debug
		rlsmgrconfig.DryRun = dryRun
		rlsmgrconfig.VerboseMode = verbose
		rlsmgrconfig.Backend = &config.BackendConfig{
			StoragePath: storagePath,
		}
		// check env for KUBECONFIG
		if kubeConfig == "" && os.Getenv(constants.EnvKubeConfig) != "" {
			kubeConfig = os.Getenv(constants.EnvKubeConfig)
		}
		rlsmgrconfig.ClusterConfig = &config.ClusterConfig{
			KubeConfig:  kubeConfig,
			KubeContext: kubeContext,
		}
		rlsmgrconfig.Helm = &config.HelmConfig{
			TillerNamespace: tillerNamespace,
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rlsmgrconfig = &config.Config{}
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Enable debugging output")
	RootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "", false, "Print planned actions without making any modifications")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	RootCmd.PersistentFlags().StringVarP(&kubeConfig, "kubeconfig", "", "", "Use this kubeconfig path, otherwise use the environment variable KUBECONFIG or ~/.kube/config")
	RootCmd.PersistentFlags().StringVarP(&kubeContext, "kubecontext", "", "", "Use this kube context, otherwise use the default")
	RootCmd.PersistentFlags().StringVarP(&storagePath, "path", "", "", "Required. Use this path within the backend for state storage")
	RootCmd.PersistentFlags().StringVarP(&tillerNamespace, "namespace", "n", "kube-system", "Communicate with the instance of Tiller in this namespace")
}
