/*
Copyright © 2020 Splice Machine
Author: Christopher Maahs <cmaahs@splicemachine.com>
*/

// Package cmd - Get Command
package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the show command
var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get various output from cloud providers",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {

	rootCmd.AddCommand(getCmd)

}
