package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/delete"
	"github.com/logicmonitor/k8s-release-manager/pkg/export"
	"github.com/logicmonitor/k8s-release-manager/pkg/healthz"
	"github.com/logicmonitor/k8s-release-manager/pkg/state"
	"github.com/logicmonitor/k8s-release-manager/pkg/transfer"
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
	Run: func(cmd *cobra.Command, args []string) {
		// Instantiate the Release Manager.
		export, err := export.New(rlsmgrconfig, mgrstate)
		if err != nil {
			log.Fatalf("Failed to create Release Manager exporter: %v", err)
		}

		if daemon {
			go export.Run() // nolint: errcheck
			// Health check.
			http.HandleFunc("/healthz", healthz.HandleFunc)
			log.Fatal(http.ListenAndServe(":8080", nil))
		} else {
			err = export.Run()
			if err != nil {
				log.Errorf("%v", err)
			}
		}
	},
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
	Run: func(cmd *cobra.Command, args []string) {
		delete, err := delete.New(rlsmgrconfig, mgrstate)
		if err != nil {
			log.Fatalf("Failed to create Release Manager deleter: %v", err)
		}

		err = delete.Run()
		if err != nil {
			log.Errorf("%v", err)
		}
	},
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
	Run: func(cmd *cobra.Command, args []string) {
		transfer, err := transfer.New(rlsmgrconfig, mgrstate)
		if err != nil {
			log.Fatalf("Failed to create Release Manager transfer: %v", err)
		}

		err = transfer.Run()
		if err != nil {
			log.Errorf("%v", err)
		}
	},
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
