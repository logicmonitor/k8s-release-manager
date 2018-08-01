package cmd

import (
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/spf13/cobra"
)

var force bool
var newStoragePath string
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
a stored release. This is by design. If you to overwrite an existing release,
you should delete it first using helm delete --purge.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		rlsmgrconfig.Helm.ReleaseTimeoutSec = int64(releaseTimeoutSec)
		rlsmgrconfig.Transfer = &config.TransferConfig{
			Force:          force,
			NewStoragePath: newStoragePath,
		}
	},
}

func init() { // nolint: dupl
	importCmd.PersistentFlags().BoolVarP(&force, "force", "", false, "Skip safety checks")
	importCmd.PersistentFlags().IntVarP(&releaseTimeoutSec, "release-timeout", "", 300, "The time, in seconds, to wait for an individual Helm release to install")
	importCmd.PersistentFlags().StringVarP(&newStoragePath, "new-path", "", "", "When installing an exported Release Manager release, update the value of --path")
	RootCmd.AddCommand(importCmd)
}
