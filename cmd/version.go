package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Express the 'version' of splice-cloud-util.",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("{\"Version\": {\"SemVer\": \"%s\", \"GitCommit\": \"%s\", \"BuildDate\": \"%s\", \"GitRef\": \"%s\"}}", semVer, gitCommit, buildDate, gitRef))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
