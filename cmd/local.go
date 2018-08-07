package cmd

import (
	"fmt"
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/state"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func localPreRun(cmd *cobra.Command) {
	localOpts := &backend.LocalOpts{}

	mgrstate = &state.State{
		Backend: &backend.Local{
			BackendConfig: rlsmgrconfig.Backend,
			Opts:          localOpts,
		},
		Config: rlsmgrconfig,
	}

	err := mgrstate.Backend.Init()
	if err != nil {
		log.Fatalf("Failed to initialize the local backend: %v", err)
	}

	err = mgrstate.Init()
	if err != nil {
		log.Fatalf("Failed to initialize state: %v", err)
	}
}

var localExportCmd = &cobra.Command{ // nolint: dupl
	Use:   "local",
	Short: "Export state using the local backend",
	Long: `Export state using the local backend
Run: ` + RootCmd.Name() + ` export --help for more information about exporting`,
	PreRun: func(cmd *cobra.Command, args []string) {
		exportCmd.PreRun(cmd, args)
		localPreRun(cmd)
	},
	Run: exportRun,
}

var localClearCmd = &cobra.Command{ // nolint: dupl
	Use:   "local",
	Short: "Clear state from the local backend",
	Long: `Clear state from the local backend
Run: ` + RootCmd.Name() + ` clear --help for more information about clearing
state`,
	PreRun: func(cmd *cobra.Command, args []string) {
		clearCmd.PreRun(cmd, args)
		localPreRun(cmd)
	},
	Run: clearRun,
}

var localImportCmd = &cobra.Command{ // nolint: dupl
	Use:   "local",
	Short: "Import state from the local backend",
	Long: `Import state from the local backend
Run: ` + RootCmd.Name() + ` import --help for more information about importing`,
	PreRun: func(cmd *cobra.Command, args []string) {
		importCmd.PreRun(cmd, args)
		localPreRun(cmd)
	},
	Run: importRun,
}

func localFlags(cmd *cobra.Command) {
	err := bindConfigFlags(cmd, map[string]string{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	localFlags(localClearCmd)
	localFlags(localExportCmd)
	localFlags(localImportCmd)
	exportCmd.AddCommand(localExportCmd)
	importCmd.AddCommand(localImportCmd)
	clearCmd.AddCommand(localClearCmd)
}
