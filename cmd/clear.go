package cmd

import (
	"github.com/spf13/cobra"
)

// managecmd represents the manage command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Deleted all configured releases and state from the backend.",
	PreRun: func(cmd *cobra.Command, args []string) {
		return
	},
}

func init() {
	RootCmd.AddCommand(clearCmd)
}
