package cmd

import (
	"github.com/logicmonitor/k8s-release-manager/pkg/utilities"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsPath string

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate the Release Manager documentation",
	Run: func(cmd *cobra.Command, args []string) {
		err := utilities.EnsureDirectory(docsPath)
		if err != nil {
			log.Warnf("%v", err)
		}
		err = doc.GenMarkdownTree(RootCmd, docsPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	docsCmd.PersistentFlags().StringVarP(&docsPath, "output", "", "./docs", "Write docs to this directory")
	RootCmd.AddCommand(docsCmd)
}
