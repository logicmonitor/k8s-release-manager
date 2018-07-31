package cmd

import (
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:    "clear",
	Short:  "Deleted all configured releases and state from the backend.",
	PreRun: func(cmd *cobra.Command, args []string) {},
}

func init() {
	RootCmd.AddCommand(clearCmd)
}
