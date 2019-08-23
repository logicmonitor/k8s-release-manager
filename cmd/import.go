package cmd

import (
	"fmt"
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/logicmonitor/k8s-release-manager/pkg/importt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var force bool
var newStoragePath, namespace, target string
var releaseTimeoutSec int

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
		valid := validateCommonConfig() && validateImportConfig()
		if !valid {
			failAuth(cmd)
		}

		rlsmgrconfig.Helm.ReleaseTimeoutSec = int64(releaseTimeoutSec)
		rlsmgrconfig.Import = &config.ImportConfig{
			Force:          force,
			NewStoragePath: newStoragePath,
			Namespace: 			namespace,
			Target:					target,
		}
	},
}

func init() { // nolint: dupl
	importCmd.PersistentFlags().BoolVarP(&force, "force", "", false, "Skip safety checks")
	importCmd.PersistentFlags().IntVarP(&releaseTimeoutSec, "release-timeout", "", 300, "The time, in seconds, to wait for an individual Helm release to install")
	importCmd.PersistentFlags().StringVarP(&newStoragePath, "new-path", "", "", "When installing an exported Release Manager release, update the value of --path")
	importCmd.PersistentFlags().StringVarP(&namespace, "namespace", "", "", "Specify a specific namespace to import releases from. Required if 'target-namespace' specified")
	importCmd.PersistentFlags().StringVarP(&target, "target-namespace", "", "", "Specify a new namespace to import releases to")
	err := bindConfigFlags(importCmd, map[string]string{
		"force":          "force",
		"releaseTimeout": "polling-timeout",
		"newPath":        "new-path",
		"namespace":			"namespace",
		"target":					"target-namespace",
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