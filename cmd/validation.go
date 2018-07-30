package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func failAuth(cmd *cobra.Command) {
	err := cmd.Help()
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(0)
}

func validateConfig() bool {
	valid := true
	if storagePath == "" {
		fmt.Println("You must specify --path")
		valid = false
	}
	return valid
}

func validateS3Config() bool {
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

func validateS3Auth() bool {
	if !validateS3SessionToken() || !validateS3Tokens() {
		return false
	}
	return true
}

func validateS3Tokens() bool {
	if (accessKeyID != "" && secretAccessKey == "") || (secretAccessKey != "" && accessKeyID == "") {
		fmt.Println("You must specify both --accessKeyID and --secretAccessKey or neither")
		return false
	}
	return true
}

func validateS3SessionToken() bool {
	if sessionToken != "" && (accessKeyID == "" || secretAccessKey == "") {
		fmt.Println("If --sessionToken is specified, you must specify --accessKeyID and --secretAccessKey")
		return false
	}
	return true
}
