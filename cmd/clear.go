package cmd

import (
	"github.com/logicmonitor/k8s-release-manager/pkg/delete"
	log "github.com/sirupsen/logrus"
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

func clearRun(cmd *cobra.Command, args []string) { //nolint: dupl
	delete, err := delete.New(rlsmgrconfig, mgrstate)
	if err != nil {
		log.Fatalf("Failed to create Release Manager deleter: %v", err)
	}

	err = delete.Run()
	if err != nil {
		log.Errorf("%v", err)
	}
}
