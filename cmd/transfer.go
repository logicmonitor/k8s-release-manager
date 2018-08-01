package cmd

import (
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/spf13/cobra"
)

var force bool
var newStoragePath string
var releaseTimeoutSec int

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Deploy the releases stored in the backend to the current cluster.",
	PreRun: func(cmd *cobra.Command, args []string) {
		rlsmgrconfig.Helm.ReleaseTimeoutSec = int64(releaseTimeoutSec)
		rlsmgrconfig.Transfer = &config.TransferConfig{
			Force:          force,
			NewStoragePath: newStoragePath,
		}
	},
}

func init() { // nolint: dupl
	transferCmd.PersistentFlags().BoolVarP(&force, "force", "", false, "Skip state safety checks.")
	transferCmd.PersistentFlags().IntVarP(&releaseTimeoutSec, "release-timeout", "", 300, "Time in seconds to wait for an individual release to install.")
	transferCmd.PersistentFlags().StringVarP(&newStoragePath, "new-path", "", "", "Updated path that the transferred release manager release will use for storing releases in the backend.")
	RootCmd.AddCommand(transferCmd)
}
