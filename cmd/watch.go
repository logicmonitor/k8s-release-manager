package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var managerReleaseName string
var pollingInterval int
var storagePath string
var tillerNamespace string

// managecmd represents the manage command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Periodically poll Tiller and export releases",
	PreRun: func(cmd *cobra.Command, args []string) {
		rlsmgrconfig.ManagerReleaseName = managerReleaseName
		rlsmgrconfig.PollingInterval = int64(pollingInterval)
		rlsmgrconfig.StoragePath = storagePath
		rlsmgrconfig.TillerNamespace = tillerNamespace
	},
}

func init() {
	log.SetLevel(log.DebugLevel)
	watchCmd.PersistentFlags().IntVarP(&pollingInterval, "polling-interval", "p", 30, "Frequency in seconds for exporting releases.")
	watchCmd.PersistentFlags().StringVarP(&managerReleaseName, "manager-release", "", "", "The Helm release name used to install the Release Manager. This value is used to prevent configuration conflicts when restoring a saved state to a new cluster.")
	watchCmd.PersistentFlags().StringVarP(&storagePath, "path", "", "", "Required. Path for storing releases in the storage backend.")
	watchCmd.PersistentFlags().StringVarP(&tillerNamespace, "namespace", "n", "kube-system", "The namespace where Tiller is deployed.")
	RootCmd.AddCommand(watchCmd)
}
