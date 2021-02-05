// Package cmd - Get Environments Code
package cmd

import (
	"strings"

	"splice-cloud-util/cmd/objects"
	"splice-cloud-util/cmd/provider"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getCmdEnvironment - List Active Environments
var getCmdEnvironment = &cobra.Command{
	Use:     "environment",
	Aliases: []string{"environments", "env", "envs"},
	Short:   "Get a list of running K8s environments for a cloud provider",
	Long: `EXAMPLE 1
	#> splice-cloud-util get environment --csp aws


EXAMPLE 2

	#> splice-cloud-util get environment --csp azure



EXAMPLE 3

	#> splice-cloud-util get environment --csp gcp

`,
	Run: func(cmd *cobra.Command, args []string) {

		csp, _ := cmd.Flags().GetString("csp")
		env, _ := cmd.Flags().GetString("environment")

		err := getActiveEnvironment(csp, env, outputFormat)
		if err != nil {
			logrus.WithError(err).Error("Error getting the active environments")
		}

	},
}

func getActiveEnvironment(csp string, env string, format string) error {

	var out objects.DetailList

	cspProvider, err := provider.GetProvider(csp, jenkinsUser, jenkinsKey, jenkinsURL)
	if err != nil {
		logrus.Fatal("Failed to get provider")
	}

	out, _ = cspProvider.GetEnvironments(env)

	switch strings.ToLower(format) {
	case "json":
		out.ToJSON()
	case "yaml":
		out.ToYAML()
	case "text", "table":
		out.ToTEXT(noHeaders)
	}

	return nil

}

func init() {

	getCmdEnvironment.Flags().StringP("csp", "c", "gcp", "The cloud provider from which to list active K8s environments")
	getCmdEnvironment.Flags().StringP("environment", "e", "", "The environment name (nonprod-eks-dev1, nonprod-az-dev1, nonprod-gke-dev1")

	getCmd.AddCommand(getCmdEnvironment)

}
