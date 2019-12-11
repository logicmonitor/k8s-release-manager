package cmd

import (
	"fmt"
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/importt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var force bool
var newStoragePath, namespace, target string
var releaseTimeoutSec, threads int
var values map[string]string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Helm release state",
	Long: `Release Manager Import will retrieve state from the configured
backend and install all exported releases to the current Kubernetes cluster.

We should avoid introducing scenarios where an imported Release Manager release
is configured to write to the same backend path as an existing
Release Manager in a different cluster. If a Release Manager release is stored
in the remote state, and --new-path is not set, this command will fail. If
you're really sure that this is an operation you want to perform (it probably
isn't), you can set --force to ignore safety checks.

Import is designed to fail if a release already exists with the same name as
a stored release. This is by design. If you want to overwrite an existing
release, you should use the helm delete --purge to delete it first.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		valid := validateCommonConfig()
		if !valid {
			failAuth(cmd)
		}

		_ = viper.GetStringMapString("valueUpdates")
		rlsmgrconfig.Helm.ReleaseTimeoutSec = viper.GetInt64("releaseTimeout")
		rlsmgrconfig.Import = &config.ImportConfig{
			Force:             viper.GetBool("force"),
			NewStoragePath:    viper.GetString("newPath"),
			Namespace:         viper.GetString("namespace"),
			Target:            viper.GetString("target"),
			Values:            values,
			ExcludeNamespaces: viper.GetStringSlice("excludeNamespaces"),
			Threads:           viper.GetInt64("threads"),
		}

		valid = validateImportConfig()
		if !valid {
			failAuth(cmd)
		}
	},
}

func init() { // nolint: dupl
	values = map[string]string{}
	importCmd.PersistentFlags().BoolVarP(&force, "force", "", false, "Skip safety checks")
	importCmd.PersistentFlags().IntVarP(&releaseTimeoutSec, "release-timeout", "", 300, "The time, in seconds, to wait for an individual Helm release to install")
	importCmd.PersistentFlags().StringVarP(&newStoragePath, "new-path", "", "", "When installing an exported Release Manager release, update the value of --path")
	importCmd.PersistentFlags().StringVarP(&namespace, "namespace", "", "", "Specify a specific namespace to import releases from. Required if 'target-namespace' specified")
	importCmd.PersistentFlags().StringVarP(&target, "target-namespace", "", "", "Specify a new namespace to import releases to")
	importCmd.PersistentFlags().StringToStringVarP(&values, "update-values", "", map[string]string{}, "Specify a mapping of values to update when importing releases. Overrides apply to all releases for which a given value is already set, but will not insert the value if it doesn't already exist")
	importCmd.PersistentFlags().StringSliceP("exclude-namespaces", "", []string{}, "A list of namespaces to exclude. The default behavior is to import all namespaces")
	importCmd.PersistentFlags().IntVarP(&threads, "threads", "", 50, "The maximum number of threads to use for installing releases")

	err := bindConfigFlags(importCmd, map[string]string{
		"force":             "force",
		"releaseTimeout":    "release-timeout",
		"newPath":           "new-path",
		"namespace":         "namespace",
		"target":            "target-namespace",
		"valueUpdates":      "update-values",
		"excludeNamespaces": "exclude-namespaces",
		"threads":           "threads",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	RootCmd.AddCommand(importCmd)
}

func importRun(cmd *cobra.Command, args []string) { // nolint: dupl
	importt, err := importt.New(rlsmgrconfig, mgrstate)
	if err != nil {
		log.Fatalf("Failed to create Release Manager import: %v", err)
	}

	err = importt.Run()
	if err != nil {
		log.Errorf("%v", err)
	}
}
