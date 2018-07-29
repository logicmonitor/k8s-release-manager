package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/logicmonitor/k8s-release-manager/pkg/healthz"
	"github.com/logicmonitor/k8s-release-manager/pkg/manager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var accessKeyID string
var bucket string
var region string
var secretAccessKey string
var sessionToken string

// managecmd represents the manage command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Periodically poll Tiller and export releases to an S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		watchCmd.PreRun(cmd, args)
		validate(cmd)

		b := &backend.S3{
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

		// Instantiate the Release Manager.
		rlsmgr, err := manager.New(rlsmgrconfig, b)
		if err != nil {
			log.Fatalf("Failed to create Release Manager: %v", err)
		}

		// Start the Release Manager.
		go rlsmgr.Run() // nolint: errcheck

		// Health check.
		http.HandleFunc("/healthz", healthz.HandleFunc)
		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

func init() {
	s3Cmd.PersistentFlags().StringVarP(&accessKeyID, "accessKeyID", "", "", "AWS access key ID with permissions to access to configured S3 bucket.")
	s3Cmd.PersistentFlags().StringVarP(&bucket, "bucket", "", "", "The S3 bucket for storing exported releases.")
	s3Cmd.PersistentFlags().StringVarP(&region, "region", "", "", "The S3 bucket region.")
	s3Cmd.PersistentFlags().StringVarP(&secretAccessKey, "secretAccessKey", "", "", "AWS secret access key with permissions to access to configured S3 bucket.")
	s3Cmd.PersistentFlags().StringVarP(&sessionToken, "sessionToken", "", "", "AWS STS session token.")
	watchCmd.AddCommand(s3Cmd)
}

func validate(cmd *cobra.Command) {
	valid := validateConfig() && validateAuth() && validateS3()

	if !valid {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err.Error())
		}
		os.Exit(0)
	}
}

func validateConfig() bool {
	valid := true
	if storagePath == "" {
		fmt.Println("You must specify --path")
		valid = false
	}
	return valid
}

func validateS3() bool {
	valid := true
	if bucket == "" {
		fmt.Println("You must specify --bucket")
		valid = false
	}
	if region == "" {
		fmt.Println("You must specify --region")
		valid = false
	}
	return valid
}

func validateAuth() bool {
	valid := true
	if (accessKeyID != "" && secretAccessKey == "") || (secretAccessKey != "" && accessKeyID == "") {
		fmt.Println("You must specify both --accessKeyID and --secretAccessKey or neither")
		valid = false
	}

	if sessionToken != "" && (accessKeyID == "" || secretAccessKey == "") {
		fmt.Println("If --sessionToken is specified, you must specify --accessKeyID and --secretAccessKey")
		valid = false
	}
	return valid
}
