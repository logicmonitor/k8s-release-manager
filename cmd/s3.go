package cmd

import (
	"net/http"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/delete"
	"github.com/logicmonitor/k8s-release-manager/pkg/export"
	"github.com/logicmonitor/k8s-release-manager/pkg/healthz"
	"github.com/logicmonitor/k8s-release-manager/pkg/state"
	"github.com/logicmonitor/k8s-release-manager/pkg/transfer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var accessKeyID string
var bucket string
var mgrstate *state.State
var region string
var s3Backend backend.Backend
var secretAccessKey string
var sessionToken string

func s3PreRun(cmd *cobra.Command) {
	valid := validateS3Auth() && validateS3Config()
	if !valid {
		failAuth(cmd)
	}

	s3Backend = &backend.S3{
		BackendConfig: rlsmgrconfig.Backend,
		Opts: &backend.S3Opts{
			Auth: &backend.S3Auth{
				AccessKeyID:     accessKeyID,
				SecretAccessKey: secretAccessKey,
				SessionToken:    sessionToken,
			},
			Bucket: bucket,
			Region: region,
		},
	}
	mgrstate = &state.State{
		Backend: s3Backend,
		Config:  rlsmgrconfig,
	}
	err := mgrstate.Init()
	if err != nil {
		log.Fatalf("Failed to initialize state: %v", err)
	}
}

var s3ExportCmd = &cobra.Command{
	Use:   "s3",
	Short: "Export state using the S3 backend",
	Long: `Export state using the S3 backend
Run: ` + RootCmd.Name() + ` export --help for more information about exporting`,
	PreRun: func(cmd *cobra.Command, args []string) {
		exportCmd.PreRun(cmd, args)
		s3PreRun(cmd)
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

var s3ClearCmd = &cobra.Command{ // nolint: dupl
	Use:   "s3",
	Short: "Clear state from the S3 backend",
	Long: `Clear state from the S3 backend
Run: ` + RootCmd.Name() + ` clear --help for more information about clearing
state`,
	PreRun: func(cmd *cobra.Command, args []string) {
		clearCmd.PreRun(cmd, args)
		s3PreRun(cmd)
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

var s3importCmd = &cobra.Command{ // nolint: dupl
	Use:   "s3",
	Short: "Import state from the S3 backend",
	Long: `Import state from the S3 backend
Run: ` + RootCmd.Name() + ` import --help for more information about importing`,
	PreRun: func(cmd *cobra.Command, args []string) {
		importCmd.PreRun(cmd, args)
		s3PreRun(cmd)
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

func s3Flags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&accessKeyID, "accessKeyID", "", "", "An AWS Access Key ID for accessing the S3 bucket, otherwise use the default AWS credential provider chain")
	cmd.PersistentFlags().StringVarP(&bucket, "bucket", "", "", "Required. Use this S3 bucket for backend storage")
	cmd.PersistentFlags().StringVarP(&region, "region", "", "us-east-1", "The backend S3 bucket's region")
	cmd.PersistentFlags().StringVarP(&secretAccessKey, "secretAccessKey", "", "", "An AWS Secret Access Key for accessing the S3 bucket, otherwise use the default AWS credential provider chain")
	cmd.PersistentFlags().StringVarP(&sessionToken, "sessionToken", "", "", "An AWS STS Session Token  for accessing the S3 bucket, otherwise use the default AWS credential provider chain")
}

func init() {
	s3Flags(s3ClearCmd)
	s3Flags(s3ExportCmd)
	s3Flags(s3importCmd)
	exportCmd.AddCommand(s3ExportCmd)
	importCmd.AddCommand(s3importCmd)
	clearCmd.AddCommand(s3ClearCmd)
}
