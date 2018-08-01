package cmd

import (
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:    "clear",
	Short:  "Clear all state",
	PreRun: func(cmd *cobra.Command, args []string) {},
}

func init() {
	RootCmd.AddCommand(clearCmd)
}
