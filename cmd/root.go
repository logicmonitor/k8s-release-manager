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
	Use: "k8s-release-manager",
	Short: `
Release Manager is a tool for periodically exporting Helm releases to
an external storage location`,
	Long: `
Release Manager provides a long-running application that will periodically
poll the Tiller server installed in your cluster. The application will
retrieve all installed releases in the cluster and write them to the
configured external storage location. The exported data can be retrieved in
the future and used the deploy the exported Helm releases to a different
cluster. The intended use case is to simplify cluster replication actions such
blue/green cluster deployments and disaster recovery scenarios.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
		if debug {
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
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Debug output.")
	RootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "", false, "Don't modify any state in the remote backend.")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output. Generally, this will print chart contents to stdout before performing an operation.")
	RootCmd.PersistentFlags().StringVarP(&kubeConfig, "kubeconfig", "", "", "Path to local kubecontext for interfacing with Helm.")
	RootCmd.PersistentFlags().StringVarP(&kubeContext, "kubecontext", "", "", "Kube context for interfacing with Helm.")
	RootCmd.PersistentFlags().StringVarP(&storagePath, "path", "", "", "Required. Path for storing releases in the storage backend.")
	RootCmd.PersistentFlags().StringVarP(&tillerNamespace, "namespace", "n", "kube-system", "The namespace where Tiller is deployed.")
}
