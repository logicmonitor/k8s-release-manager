package cmd

import (
	"net/http"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/delete"
	"github.com/logicmonitor/k8s-release-manager/pkg/export"
	"github.com/logicmonitor/k8s-release-manager/pkg/healthz"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var accessKeyID string
var bucket string
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
}

var s3ExportCmd = &cobra.Command{
	Use:   "s3",
	Short: "Use the s3 backend",
	PreRun: func(cmd *cobra.Command, args []string) {
		exportCmd.PreRun(cmd, args)
		s3PreRun(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Instantiate the Release Manager.
		export, err := export.New(rlsmgrconfig, s3Backend)
		if err != nil {
			log.Fatalf("Failed to create Release Manager exporter: %v", err)
		}

		// Start the exporter.
		if daemon {
			go export.Run() // nolint: errcheck
			// Health check.
			http.HandleFunc("/healthz", healthz.HandleFunc)
			log.Fatal(http.ListenAndServe(":8080", nil))
		} else {
			export.Run()
		}
	},
}

var s3ClearCmd = &cobra.Command{
	Use:   "s3",
	Short: "Use the s3 backend",
	PreRun: func(cmd *cobra.Command, args []string) {
		clearCmd.PreRun(cmd, args)
		s3PreRun(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		delete, err := delete.New(rlsmgrconfig, s3Backend)
		if err != nil {
			log.Fatalf("Failed to create Release Manager deleter: %v", err)
		}

		// Start the deleter.
		delete.Run() // nolint: errcheck
	},
}

func s3Flags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&accessKeyID, "accessKeyID", "", "", "AWS access key ID with permissions to access to configured S3 bucket.")
	cmd.PersistentFlags().StringVarP(&bucket, "bucket", "", "", "The S3 bucket for storing exported releases.")
	cmd.PersistentFlags().StringVarP(&region, "region", "", "", "The S3 bucket region.")
	cmd.PersistentFlags().StringVarP(&secretAccessKey, "secretAccessKey", "", "", "AWS secret access key with permissions to access to configured S3 bucket.")
	cmd.PersistentFlags().StringVarP(&sessionToken, "sessionToken", "", "", "AWS STS session token.")
}

func init() {
	s3Flags(s3ExportCmd)
	s3Flags(s3ClearCmd)
	exportCmd.AddCommand(s3ExportCmd)
	clearCmd.AddCommand(s3ClearCmd)
}
