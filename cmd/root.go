package cmd

import (
	"fmt"
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/constants"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rlsmgrconfig *config.Config
var cfgFile string
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
	Use:   "releasemanager",
	Short: "Release Manager is a tool for importing and exporting Helm release state",
	Long: `
Release Manager provides functionality for exporting and importing the
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
		rlsmgrconfig.DebugMode = viper.GetBool("debug")
		rlsmgrconfig.DryRun = viper.GetBool("dryRun")
		rlsmgrconfig.VerboseMode = viper.GetBool("verbose")
		rlsmgrconfig.Backend = &config.BackendConfig{
			StoragePath: viper.GetString("path"),
		}

		// check env for KUBECONFIG
		kubeConfig = viper.GetString("kubeconfig")
		if kubeConfig == "" && os.Getenv(constants.EnvKubeConfig) != "" {
			kubeConfig = os.Getenv(constants.EnvKubeConfig)
		}

		rlsmgrconfig.ClusterConfig = &config.ClusterConfig{
			KubeConfig:  kubeConfig,
			KubeContext: viper.GetString("kubecontext"),
		}

		rlsmgrconfig.Helm = &config.HelmConfig{
			TillerNamespace: viper.GetString("tillerNamespace"),
		}
		if rlsmgrconfig.DebugMode {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
		if rlsmgrconfig.DryRun {
			fmt.Println("Dry run. No changes will be made.")
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
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Set a custom configuration path")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Enable debugging output")
	RootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "", false, "Print planned actions without making any modifications")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	RootCmd.PersistentFlags().StringVarP(&kubeConfig, "kubeconfig", "", "", "Use this kubeconfig path, otherwise use the environment variable KUBECONFIG or ~/.kube/config")
	RootCmd.PersistentFlags().StringVarP(&kubeContext, "kubecontext", "", "", "Use this kube context, otherwise use the default")
	RootCmd.PersistentFlags().StringVarP(&storagePath, "path", "", "", "Required. Use this path within the backend for state storage")
	RootCmd.PersistentFlags().StringVarP(&tillerNamespace, "tiller-namespace", "", "kube-system", "Communicate with the instance of Tiller in this namespace")
	err := bindConfigFlags(RootCmd, map[string]string{
		"debug":           "debug",
		"dryRun":          "dry-run",
		"verbose":         "verbose",
		"kubeconfig":      "kubeconfig",
		"kubecontext":     "kubecontext",
		"path":            "path",
		"tillerNamespace": "tiller-namespace",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func bindConfigFlags(cmd *cobra.Command, mapping map[string]string) (err error) {
	for k, v := range mapping {
		err = viper.BindPFlag(k, cmd.PersistentFlags().Lookup(v))
		if err != nil {
			return err
		}
	}
	return err
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(constants.DefaultConfigPath)
		viper.SetConfigName("config")
	}

	err := viper.ReadInConfig()
	if err != nil && cfgFile != "" {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
