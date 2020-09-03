package cmd

import (
	"fmt"
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/state"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accessKeyID string
var bucket string
var region string
var secretAccessKey string
var sessionToken string

func s3PreRun(cmd *cobra.Command) {
	s3Opts := &backend.S3Opts{
		Auth: &backend.S3Auth{
			AccessKeyID:     viper.GetString("accessKeyID"),
			SecretAccessKey: viper.GetString("secretAccessKey"),
			SessionToken:    viper.GetString("sessionToken"),
		},
		Bucket: viper.GetString("bucket"),
		Region: viper.GetString("region"),
	}

	valid := validateS3Auth(s3Opts) && validateS3Config(s3Opts)
	if !valid {
		failAuth(cmd)
	}

	mgrstate = &state.State{
		Backend: &backend.S3{
			BackendConfig: rlsmgrconfig.Backend,
			Opts:          s3Opts,
		},
		Config: rlsmgrconfig,
	}
	err := mgrstate.Init()
	if err != nil {
		log.Fatalf("Failed to initialize state: %v", err)
	}

}

var s3ExportCmd = &cobra.Command{ // nolint: dupl
	Use:   "s3",
	Short: "Export state using the S3 backend",
	Long: `Export state using the S3 backend
Run: ` + RootCmd.Name() + ` export --help for more information about exporting`,
	PreRun: func(cmd *cobra.Command, args []string) {
		exportCmd.PreRun(cmd, args)
		s3PreRun(cmd)
	},
	Run: exportRun,
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
	Run: clearRun,
}

var s3ImportCmd = &cobra.Command{ // nolint: dupl
	Use:   "s3",
	Short: "Import state from the S3 backend",
	Long: `Import state from the S3 backend
Run: ` + RootCmd.Name() + ` import --help for more information about importing`,
	PreRun: func(cmd *cobra.Command, args []string) {
		importCmd.PreRun(cmd, args)
		s3PreRun(cmd)
	},
	Run: importRun,
}

func s3Flags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&accessKeyID, "accessKeyID", "", "", "An AWS Access Key ID for accessing the S3 bucket, otherwise use the default AWS credential provider chain")
	cmd.PersistentFlags().StringVarP(&bucket, "bucket", "", "", "Required. Use this S3 bucket for backend storage")
	cmd.PersistentFlags().StringVarP(&region, "region", "", "us-east-1", "The backend S3 bucket's region")
	cmd.PersistentFlags().StringVarP(&secretAccessKey, "secretAccessKey", "", "", "An AWS Secret Access Key for accessing the S3 bucket, otherwise use the default AWS credential provider chain")
	cmd.PersistentFlags().StringVarP(&sessionToken, "sessionToken", "", "", "An AWS STS Session Token  for accessing the S3 bucket, otherwise use the default AWS credential provider chain")
	err := bindConfigFlags(cmd, map[string]string{
		"accessKeyID":     "accessKeyID",
		"bucket":          "bucket",
		"region":          "region",
		"secretAccessKey": "secretAccessKey",
		"sessionToken":    "sessionToken",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	s3Flags(s3ClearCmd)
	s3Flags(s3ExportCmd)
	s3Flags(s3ImportCmd)
	exportCmd.AddCommand(s3ExportCmd)
	importCmd.AddCommand(s3ImportCmd)
	clearCmd.AddCommand(s3ClearCmd)
}
