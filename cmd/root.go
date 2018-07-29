package cmd

import (
	"fmt"
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rlsmgrconfig *config.Config
var cfgFile string
var debug bool
var verbose bool

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
		// Retrieve the application configuration.
		var err error
		rlsmgrconfig, err = config.New()
		if err != nil {
			log.Fatalf("Failed to get config: %v", err)
		}
		rlsmgrconfig.DebugMode = debug
		rlsmgrconfig.VerboseMode = verbose
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
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-release-manager.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Debug mode. Log release info and skip writing to backend.")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		setDefaultConfig()
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setDefaultConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".k8s-release-manager" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".k8s-release-manager")
}
