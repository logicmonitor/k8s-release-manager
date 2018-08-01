package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsPath string

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate the Release Manager documentation",
	Run: func(cmd *cobra.Command, args []string) {
		os.MkdirAll(docsPath, os.ModePerm)
		err := doc.GenMarkdownTree(RootCmd, docsPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	docsCmd.PersistentFlags().StringVarP(&docsPath, "output", "", "./docs", "Write docs to this directory")
	RootCmd.AddCommand(docsCmd)
}
