package cmd

import (
	"fmt"
	"os"

	"github.com/logicmonitor/k8s-release-manager/pkg/backend"
	"github.com/spf13/cobra"
)

func failAuth(cmd *cobra.Command) {
	err := cmd.Help()
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(0)
}

func validateCommonConfig() bool {
	valid := true
	if rlsmgrconfig.Backend.StoragePath == "" {
		fmt.Println("You must specify --path")
		valid = false
	}
	return valid
}

func validateS3Config(opts *backend.S3Opts) bool {
	valid := true
	if opts.Bucket == "" {
		fmt.Println("You must specify --bucket")
		valid = false
	}
	if opts.Region == "" {
		fmt.Println("You must specify --region")
		valid = false
	}
	return valid
}

func validateS3Auth(opts *backend.S3Opts) bool {
	if !validateS3SessionToken(opts) || !validateS3Tokens(opts) {
		return false
	}
	return true
}

func validateS3Tokens(opts *backend.S3Opts) bool {
	if (opts.Auth.AccessKeyID != "" && opts.Auth.SecretAccessKey == "") || (opts.Auth.AccessKeyID == "" && opts.Auth.SecretAccessKey != "") {
		fmt.Println("You must specify both --accessKeyID and --secretAccessKey or neither")
		return false
	}
	return true
}

func validateS3SessionToken(opts *backend.S3Opts) bool {
	if opts.Auth.SessionToken != "" && (opts.Auth.AccessKeyID == "" || opts.Auth.SecretAccessKey == "") {
		fmt.Println("If --sessionToken is specified, you must specify --accessKeyID and --secretAccessKey")
		return false
	}
	return true
}
