package cmd

import (
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
	"github.com/spf13/cobra"
)

var daemon bool
var pollingInterval int

// managecmd represents the manage command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export installed Helm releases to the configured backend.",
	PreRun: func(cmd *cobra.Command, args []string) {
		valid := validateConfig()
		if !valid {
			failAuth(cmd)
		}

		rlsmgrconfig.Export = &config.ExportConfig{
			DaemonMode:      daemon,
			ReleaseName:     releaseName,
			PollingInterval: int64(pollingInterval),
		}
		rlsmgrconfig.Helm = &config.HelmConfig{
			TillerNamespace: tillerNamespace,
			TillerHost:      tillerHost,
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "", false, "Periodically poll Tiller and export releases.")
	exportCmd.PersistentFlags().IntVarP(&pollingInterval, "polling-interval", "p", 30, "Frequency in seconds for exporting releases.")
	exportCmd.PersistentFlags().StringVarP(&releaseName, "release-name", "", "", "The Helm release name used to install the Release Manager. This value is used to prevent configuration conflicts when restoring a saved state to a new cluster.")
	RootCmd.AddCommand(exportCmd)
}
