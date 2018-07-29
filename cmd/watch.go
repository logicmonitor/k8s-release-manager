package cmd

import (
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pollingInterval int
var releaseName string
var storagePath string
var tillerHost string
var tillerNamespace string

// managecmd represents the manage command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Periodically poll Tiller and export releases",
	PreRun: func(cmd *cobra.Command, args []string) {
		rlsmgrconfig.Manager = &config.ManagerConfig{
			ReleaseName:     releaseName,
			PollingInterval: int64(pollingInterval),
			StoragePath:     storagePath,
		}
		rlsmgrconfig.Helm = &config.HelmConfig{
			TillerNamespace: tillerNamespace,
			TillerHost:      tillerHost,
		}
	},
}

func init() {
	log.SetLevel(log.DebugLevel)
	watchCmd.PersistentFlags().IntVarP(&pollingInterval, "polling-interval", "p", 30, "Frequency in seconds for exporting releases.")
	watchCmd.PersistentFlags().StringVarP(&releaseName, "manager-release", "", "", "The Helm release name used to install the Release Manager. This value is used to prevent configuration conflicts when restoring a saved state to a new cluster.")
	watchCmd.PersistentFlags().StringVarP(&storagePath, "path", "", "", "Required. Path for storing releases in the storage backend.")
	watchCmd.PersistentFlags().StringVarP(&tillerHost, "tiller-host", "n", "kube-system", "The namespace where Tiller is deployed.")
	watchCmd.PersistentFlags().StringVarP(&tillerNamespace, "namespace", "n", "kube-system", "The namespace where Tiller is deployed.")
	RootCmd.AddCommand(watchCmd)
}
